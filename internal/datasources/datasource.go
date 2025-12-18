package datasources

import (
	"context"

	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/core/models"
)

type Datasource interface {
	AddUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
	AddPassword(ctx context.Context, login string, password *models.Password) error
	GetPassword(ctx context.Context, login string, passwordID int) (*models.Password, error)
	AddSecret(ctx context.Context, userID models.UserID, secretType models.SecretType, secret []byte, comment string) error
	GetSecret(ctx context.Context, secretID int, userID models.UserID) (*models.Secret, error)
	GetAllSecrets(ctx context.Context, userID models.UserID) ([]models.SecretInfo, error)
}

type dummyDS struct {
	cfg config.ConfigProvider
}

func NewDummyDS(cfg config.ConfigProvider) Datasource {
	return &dummyDS{cfg: cfg}
}

func (d dummyDS) AddUser(ctx context.Context, user *models.User) (*models.User, error) {
	d.cfg.Logger().Info("RegisterUser called", zap.String("username", user.Username), zap.String("password", user.Password))

	return nil, nil
}

func (d dummyDS) GetUser(ctx context.Context, username string) (*models.User, error) {
	d.cfg.Logger().Info("GetUser called", zap.String("username", username))

	return nil, nil
}

func (d dummyDS) AddPassword(ctx context.Context, login string, password *models.Password) error {
	d.cfg.Logger().Info("AddPassword called", zap.String("login", login), zap.String("password", password.Password))

	return nil
}

func (d dummyDS) GetPassword(ctx context.Context, login string, passwordID int) (*models.Password, error) {
	d.cfg.Logger().Info("AddPassword called", zap.String("login", login), zap.Int("passwordID", passwordID))

	return nil, nil
}

func (d dummyDS) AddSecret(ctx context.Context, userID models.UserID, secretType models.SecretType, secret []byte, comment string) error {
	d.cfg.Logger().Info("AddSecret called", zap.String("secretType", string(secretType)))
	return nil
}

func (d dummyDS) GetSecret(ctx context.Context, secretID int, userID models.UserID) (*models.Secret, error) {
	d.cfg.Logger().Info("GetSecret called", zap.Int("secretID", secretID))
	return nil, nil
}

func (d dummyDS) GetAllSecrets(ctx context.Context, userID models.UserID) ([]models.SecretInfo, error) {
	d.cfg.Logger().Info("GetAllSecrets called", zap.Int("userID", int(userID)))
	return nil, nil
}
