package middleware

import (
	"gophkeep/internal/config"
	"gophkeep/internal/domain/service"
)

type MiddlewareProvider interface {
	BearerValidator
	SizeChecker
}
type MiddlewareProviderImpl struct {
	BearerValidator
	SizeChecker
}

func NewMiddlewareProvider(cfg config.ConfigProvider, serviceProvider service.IService) MiddlewareProvider {
	return &MiddlewareProviderImpl{
		BearerValidator: NewBearerValidator(cfg, serviceProvider),
		SizeChecker:     NewSizeChecker(cfg),
	}
}
