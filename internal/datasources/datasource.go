package datasources

import (
	"context"

	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/domain/models"
)

type Datasource interface {
	AddUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, username string) (*models.User, error)
	AddPassword(ctx context.Context, login string, password *models.Password) error
	GetPassword(ctx context.Context, login string, passwordID int) (*models.Password, error)
}

type dummyDS struct {
	cfg config.ConfigProvider
}

func NewDummyDS(cfg config.ConfigProvider) Datasource {
	return &dummyDS{cfg: cfg}
}

func (d dummyDS) AddUser(ctx context.Context, user *models.User) error {
	d.cfg.Logger().Info("AddUser called", zap.String("username", user.Username), zap.String("password", user.Password))

	return nil
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
