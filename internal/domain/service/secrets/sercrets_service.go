package secrets

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/core/models"
	"gophkeep/internal/datasources"
)

type SecretsService interface {
	AddPassword(ctx context.Context, object *models.PasswordRequestObject) error
	AddSecret(ctx context.Context, object *models.Secret) error
	GetAllSecrets(ctx context.Context) ([]models.SecretInfo, error)
	GetSecret(ctx context.Context, secretID int) (*models.Secret, error)
}

type secretsService struct {
	db  datasources.Datasource
	cfg config.ConfigProvider
}

func NewSecretsService(cfg config.ConfigProvider, db datasources.Datasource) SecretsService {
	return &secretsService{cfg: cfg, db: db}
}

func (s secretsService) AddPassword(ctx context.Context, object *models.PasswordRequestObject) error {
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

	if err := s.db.AddSecret(ctx, userID, models.PASSWORD, secret, ""); err != nil {
		return fmt.Errorf("add secret: %w", err)
	}

	return nil
}

func (s secretsService) AddSecret(ctx context.Context, object *models.Secret) error {
	s.cfg.Logger().Info("Receive add secret request", zap.String("type", string(object.Type)), zap.String("comment", object.Comment))

	userID, ok := ctx.Value("userID").(models.UserID)
	if !ok {
		return fmt.Errorf("encode secret: userID not found in context")
	}

	if err := s.db.AddSecret(ctx, userID, object.Type, object.Content, object.Comment); err != nil {
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

	s.cfg.Logger().Info("Get secret", zap.Int("secretID", secretID))

	return secret, nil
}
