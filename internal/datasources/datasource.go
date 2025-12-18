package datasources

import (
	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/domain/models"
)

type Datasource interface {
	AddUser(username string, password string) error
	GetUser(username string) (*models.User, error)
	AddPassword(username string, password *models.Password) error
	GetPassword(username string, passwordID int) (*models.Password, error)
}

type dummyDS struct {
	cfg config.ConfigProvider
}

func NewDummyDS(cfg config.ConfigProvider) Datasource {
	return &dummyDS{cfg: cfg}
}

func (d dummyDS) AddUser(username string, password string) error {
	d.cfg.Logger().Info("AddUser called", zap.String("username", username), zap.String("password", password))

	return nil
}

func (d dummyDS) GetUser(username string) (*models.User, error) {
	d.cfg.Logger().Info("GetUser called", zap.String("username", username))

	return nil, nil
}

func (d dummyDS) AddPassword(login string, password *models.Password) error {
	d.cfg.Logger().Info("AddPassword called", zap.String("login", login), zap.String("password", password.Password))

	return nil
}

func (d dummyDS) GetPassword(login string, passwordID int) (*models.Password, error) {
	d.cfg.Logger().Info("AddPassword called", zap.String("login", login), zap.Int("passwordID", passwordID))

	return nil, nil
}
