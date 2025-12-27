package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/core/models"
	"gophkeep/internal/core/requests"
	"gophkeep/internal/core/validation"
	"gophkeep/internal/datasources"
	"gophkeep/internal/domain/cryptomanager"
)

type SecretsService interface {
	AddPassword(ctx context.Context, object *requests.PasswordRequestObject) error
	AddSecret(ctx context.Context, object *requests.SecretRequest) error
	GetAllSecrets(ctx context.Context) ([]models.SecretInfo, error)
	GetSecret(ctx context.Context, secretID int) (*models.Secret, error)
	UpdateSecret(ctx context.Context, object *requests.SecretRequest) error
	DeleteSecret(ctx context.Context, secretID int) error
	Ping(ctx context.Context) error
}

type secretsService struct {
	db              datasources.Datasource
	cfg             config.ConfigProvider
	encryptor       cryptomanager.CryptoManager
	secretValidator validation.IValidator
}

func NewSecretsService(cfg config.ConfigProvider, db datasources.Datasource) SecretsService {
	s := &secretsService{cfg: cfg, db: db, secretValidator: validation.NewValidator(cfg)}

	if key := cfg.GetConfig().Key; key != "" {
		encryptor, err := cryptomanager.NewCryptoManager([]byte(key))
		if err != nil {
			panic(err)
		}
		s.encryptor = encryptor
	}

	return s
}

func (s secretsService) AddPassword(ctx context.Context, object *requests.PasswordRequestObject) error {
	s.cfg.Logger().Info("Receive add password request", zap.String("login", object.Login), zap.String("password", object.Password))

	pass := models.Password{
		Login:    object.Login,
		Password: object.Password,
	}

	secret, err := json.Marshal(pass)
	if err != nil {
		return fmt.Errorf("encode secret: %w", err)
	}

	userID, ok := ctx.Value("userID").(models.UserID)
	if !ok {
		return fmt.Errorf("encode secret: userID not found in context")
	}

	if _, err := s.db.AddSecret(ctx, userID, models.PASSWORD, secret, "", time.Now()); err != nil {
		return fmt.Errorf("add secret: %w", err)
	}

	return nil
}

func (s secretsService) AddSecret(ctx context.Context, object *requests.SecretRequest) error {
	s.cfg.Logger().Info("Receive add secret request", zap.String("type", string(object.SecretType)), zap.String("comment", object.Comment))

	content, err := s.secretValidator.Validate(object.BodyString, object.SecretType, object.IsEncoded)
	if err != nil {
		return fmt.Errorf("secret is not valid: %w", err)
	}

	secret := &models.Secret{
		Type:    object.SecretType,
		Comment: object.Comment,
		Content: content,
	}

	userID, ok := ctx.Value("userID").(models.UserID)
	if !ok {
		return fmt.Errorf("encode secret: userID not found in context")
	}

	secretEncrypted, err := s.encryptor.Encrypt(secret.Content)
	if err != nil {
		return fmt.Errorf("encrypt secret error: %w", err)
	}

	if _, err := s.db.AddSecret(ctx, userID, secret.Type, secretEncrypted, secret.Comment, time.Now()); err != nil {
		return fmt.Errorf("add secret: %w", err)
	}

	return nil
}

func (s secretsService) GetAllSecrets(ctx context.Context) ([]models.SecretInfo, error) {
	userID := ctx.Value("userID").(models.UserID)
	secrets, err := s.db.GetAllSecrets(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get secrets: %w", err)
	}

	return secrets, nil
}

func (s secretsService) GetSecret(ctx context.Context, secretID int) (*models.Secret, error) {
	s.cfg.Logger().Info("Receive get secret request", zap.Int("secretID", secretID))
	userID := ctx.Value("userID").(models.UserID)

	secret, err := s.db.GetSecret(ctx, secretID, userID)
	if err != nil {
		return nil, fmt.Errorf("get secret: %w", err)
	}

	secretDecrypted, err := s.decrypt(secret.Content)
	if err != nil {
		return nil, fmt.Errorf("decrypt secret: %w", err)
	}
	secret.Content = secretDecrypted

	s.cfg.Logger().Info("Get secret", zap.Int("secretID", secretID))

	return secret, nil
}

func (s secretsService) UpdateSecret(ctx context.Context, object *requests.SecretRequest) error {
	s.cfg.Logger().Info("Receive update secret request", zap.String("type", string(object.SecretType)), zap.String("comment", object.Comment))

	content, err := s.secretValidator.Validate(object.BodyString, object.SecretType, object.IsEncoded)
	if err != nil {
		return fmt.Errorf("secret is not valid: %w", err)
	}

	userID, ok := ctx.Value("userID").(models.UserID)
	if !ok {
		return fmt.Errorf("encode object: userID not found in context")
	}

	secretEncrypted, err := s.encryptor.Encrypt(content)
	if err != nil {
		return fmt.Errorf("encrypt object error: %w", err)
	}

	secret := &models.Secret{
		ID:        *object.ID,
		Type:      object.SecretType,
		Comment:   object.Comment,
		Content:   secretEncrypted,
		UpdatedAt: time.Now(),
	}

	if err := s.db.UpdateSecret(ctx, secret, userID); err != nil {
		return fmt.Errorf("update secret: %w", err)
	}

	return nil
}

func (s secretsService) DeleteSecret(ctx context.Context, secretID int) error {
	s.cfg.Logger().Info("Receive delete secret request", zap.Int("secretID", secretID))

	userID, ok := ctx.Value("userID").(models.UserID)
	if !ok {
		return fmt.Errorf("encode object: userID not found in context")
	}

	if err := s.db.DeleteSecret(ctx, secretID, userID); err != nil {
		return fmt.Errorf("delete secret: %w", err)
	}

	return nil
}

func (s secretsService) encryptWithNewKey(secret []byte) (secretEncrypted []byte, encryptionKey []byte, err error) {
	encryptionKey, err = s.encryptor.GenerateKey(32)
	if err != nil {
		return nil, nil, fmt.Errorf("generate encryption key error: %w", err)
	}
	encryptor, err := cryptomanager.NewCryptoManager(encryptionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("creating encryptor manager error: %w", err)
	}
	secretEncrypted, err = encryptor.Encrypt(secret)
	if err != nil {
		return nil, nil, fmt.Errorf("encrypt secret error: %w", err)
	}

	encryptionKey, err = s.encryptor.Encrypt(encryptionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("encrypting encryption key error: %w", err)
	}

	return secretEncrypted, encryptionKey, nil
}

func (s secretsService) encrypt(secret []byte) (secretEncrypted []byte, err error) {
	if s.encryptor == nil {
		return secret, nil
	}

	secretEncrypted, err = s.encryptor.Encrypt(secret)
	if err != nil {
		return nil, fmt.Errorf("encrypt secret error: %w", err)
	}

	return secretEncrypted, nil
}

func (s secretsService) decrypt(secret []byte) (secretDecrypted []byte, err error) {
	if s.encryptor == nil {
		return secret, nil
	}

	secretDecrypted, err = s.encryptor.Decrypt(secret)
	if err != nil {
		return nil, fmt.Errorf("decrypting secret error: %w", err)
	}

	return secretDecrypted, nil
}

func (s secretsService) Ping(ctx context.Context) error {
	if err := s.db.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	return nil
}
