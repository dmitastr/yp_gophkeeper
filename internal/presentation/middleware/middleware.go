package middleware

import (
	"gophkeep/internal/config"
	"gophkeep/internal/domain/service"
)

type MiddlewareProvider interface {
	BearerValidator
}
type MiddlewareProviderImpl struct {
	BearerValidator
}

func NewMiddlewareProvider(cfg config.ConfigProvider, serviceProvider service.IService) MiddlewareProvider {
	return &MiddlewareProviderImpl{BearerValidator: NewBearerValidator(cfg, serviceProvider)}
}
