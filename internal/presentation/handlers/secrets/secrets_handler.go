package secrets

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/core/models"
	"gophkeep/internal/core/requests"
	"gophkeep/internal/domain/service"
)

type SecretsHandler interface {
	AddPassword(ctx *gin.Context)
	AddSecret(ctx *gin.Context)
	GetAllSecrets(ctx *gin.Context)
	GetSecret(ctx *gin.Context)
}

type SecretsHandlerImpl struct {
	serviceProvider service.IService
	cfg             config.ConfigProvider
}

func NewSecretsHandler(cfg config.ConfigProvider, serviceProvider service.IService) *SecretsHandlerImpl {
	return &SecretsHandlerImpl{serviceProvider: serviceProvider, cfg: cfg}
}

func (h *SecretsHandlerImpl) AddPassword(ctx *gin.Context) {
	var request models.PasswordRequestObject
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.serviceProvider.AddPassword(ctx, &request); err != nil {
		h.cfg.Logger().Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password added successfully"})
}

func (h *SecretsHandlerImpl) AddSecret(ctx *gin.Context) {
	var request requests.SecretRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch tp := request.SecretType; tp {
	case models.PASSWORD:
		var passwordRequest models.PasswordRequestObject
		if err := json.Unmarshal(request.Body, &passwordRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.cfg.Logger().Info("Add secret endpoint called", zap.String("type", string(request.SecretType)), zap.String("comment", request.Comment))

		secret := models.Secret{
			Type:    request.SecretType,
			Content: request.Body,
			Comment: request.Comment,
		}
		if err := h.serviceProvider.AddSecret(ctx, &secret); err != nil {
			h.cfg.Logger().Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "secret type not supported"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password added successfully"})
}

func (h *SecretsHandlerImpl) GetAllSecrets(ctx *gin.Context) {
	secrets, err := h.serviceProvider.GetAllSecrets(ctx)
	if err != nil {
		h.cfg.Logger().Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"secrets": secrets})
}

func (h *SecretsHandlerImpl) GetSecret(ctx *gin.Context) {
	secretIDStr := ctx.Param("secretID")
	secretID, err := strconv.Atoi(secretIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secret, err := h.serviceProvider.GetSecret(ctx, secretID)
	if err != nil {
		h.cfg.Logger().Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"secret": secret})
}
