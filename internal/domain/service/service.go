package service

import (
	"gophkeep/internal/config"
	"gophkeep/internal/datasources"
	"gophkeep/internal/domain/service/authservice"
	"gophkeep/internal/domain/service/secrets"
)

type IService interface {
	authservice.AuthService
	secrets.SecretsService
}

type ServiceImpl struct {
	secrets.SecretsService
	authservice.AuthService
	db datasources.Datasource
}

func NewService(cfg config.ConfigProvider, db datasources.Datasource) IService {
	return &ServiceImpl{
		AuthService:    authservice.NewAuthService(cfg, db),
		SecretsService: secrets.NewSecretsService(cfg, db),
		db:             db,
	}
}
