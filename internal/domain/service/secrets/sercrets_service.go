package secrets

import (
	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/datasources"
	"gophkeep/internal/domain/models"
)

type SecretsService interface {
	AddPassword(object *models.PasswordRequestObject) error
}

type secretsService struct {
	db  datasources.Datasource
	cfg config.ConfigProvider
}

func NewSecretsService(cfg config.ConfigProvider, db datasources.Datasource) SecretsService {
	return &secretsService{cfg: cfg, db: db}
}

func (s secretsService) AddPassword(object *models.PasswordRequestObject) error {
	s.cfg.Logger().Info("Receive add password request", zap.String("login", object.Login), zap.String("password", object.Password))
	return nil
}
