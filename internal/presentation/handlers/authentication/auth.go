package authentication

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/domain/service"
	"gophkeep/internal/presentation/params"
)

// AuthHandler defines the contract for authentication-related HTTP handlers.
// It exposes methods for user registration and user login operations.
type AuthHandler interface {
	// LoginUser handles user authentication requests and returns an auth token
	// if the provided credentials are valid.
	LoginUser(*gin.Context)

	// RegisterUser handles user registration requests and returns an auth token
	// for the newly created user.
	RegisterUser(*gin.Context)
}

// AuthHandlerImpl is the concrete implementation of the AuthHandler interface.
// It delegates authentication logic to the service layer and uses configuration
// services such as logging.
type AuthHandlerImpl struct {
	serviceProvider service.IService
	cfg             config.ConfigProvider
}

// NewAuthHandler creates and returns a new AuthHandler instance
func NewAuthHandler(cfg config.ConfigProvider, serviceProvider service.IService) AuthHandler {
	return &AuthHandlerImpl{serviceProvider: serviceProvider, cfg: cfg}
}

// RegisterUser processes an HTTP request to register a new user.
// It validates the request payload, invokes the service layer to register
// the user, and returns an authentication token on success.
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

// LoginUser processes an HTTP request to authenticate an existing user.
// It validates the request payload, invokes the service layer to authenticate
// the user, and returns an authentication token on success.
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
