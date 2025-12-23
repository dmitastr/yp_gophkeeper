package authentication

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/domain/service"
	"gophkeep/internal/presentation/params"
)

type AuthHandler interface {
	LoginUser(*gin.Context)
	RegisterUser(*gin.Context)
}

type AuthHandlerImpl struct {
	serviceProvider service.IService
	cfg             config.ConfigProvider
}

func NewAuthHandler(cfg config.ConfigProvider, serviceProvider service.IService) AuthHandler {
	return &AuthHandlerImpl{serviceProvider: serviceProvider, cfg: cfg}
}

func (a *AuthHandlerImpl) RegisterUser(c *gin.Context) {
	var request params.AuthRequestObject
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.serviceProvider.RegisterUser(c, request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		a.cfg.Logger().Error("error while creating auth token", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *AuthHandlerImpl) LoginUser(c *gin.Context) {
	var request params.AuthRequestObject
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.serviceProvider.LoginUser(c, request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		a.cfg.Logger().Error("error while creating auth token", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
