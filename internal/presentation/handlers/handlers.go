package handlers

import (
	"gophkeep/internal/config"
	"gophkeep/internal/domain/service"
	"gophkeep/internal/presentation/handlers/authentication"
	"gophkeep/internal/presentation/handlers/ping"
	"gophkeep/internal/presentation/handlers/secrets"
)

type HandlersProvider interface {
	authentication.AuthHandler
	ping.PingHandler
	secrets.SecretsHandler
}

type HandlersImpl struct {
	authentication.AuthHandler
	ping.PingHandler
	secrets.SecretsHandler
	cfg config.ConfigProvider
}

func NewHandlersProvider(cfg config.ConfigProvider, serviceProvider service.IService) HandlersProvider {
	return &HandlersImpl{
		AuthHandler:    authentication.NewAuthHandler(cfg, serviceProvider),
		PingHandler:    ping.NewHandler(cfg, serviceProvider),
		SecretsHandler: secrets.NewSecretsHandler(cfg, serviceProvider),
		cfg:            cfg,
	}
}
