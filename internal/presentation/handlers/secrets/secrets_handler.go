package secrets

import (
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
	UpdateSecret(ctx *gin.Context)
	DeleteSecret(ctx *gin.Context)
}

type SecretsHandlerImpl struct {
	serviceProvider service.IService
	cfg             config.ConfigProvider
}

func NewSecretsHandler(cfg config.ConfigProvider, serviceProvider service.IService) *SecretsHandlerImpl {
	return &SecretsHandlerImpl{serviceProvider: serviceProvider, cfg: cfg}
}

func (h *SecretsHandlerImpl) AddPassword(ctx *gin.Context) {
	var request requests.PasswordRequestObject
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
		h.cfg.Logger().Error("error unmarshalling request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.cfg.Logger().Info("Add secret endpoint called", zap.String("type", string(request.SecretType)), zap.String("comment", request.Comment))

	secret := models.Secret{
		Type:          request.SecretType,
		Content:       request.Body,
		ContentString: request.BodyString,
		Comment:       request.Comment,
	}
	if err := h.serviceProvider.AddSecret(ctx, &secret); err != nil {
		h.cfg.Logger().Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Secret added successfully"})
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

func (h *SecretsHandlerImpl) UpdateSecret(ctx *gin.Context) {
	secretIDStr := ctx.Param("secretID")
	secretID, err := strconv.Atoi(secretIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request requests.SecretRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.cfg.Logger().Error("error unmarshalling request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.cfg.Logger().Info("Update secret endpoint called", zap.String("type", string(request.SecretType)), zap.Int("secretID", secretID))

	secret := &models.Secret{
		ID:      secretID,
		Type:    request.SecretType,
		Content: request.Body,
		Comment: request.Comment,
	}
	if err := h.serviceProvider.UpdateSecret(ctx, secret); err != nil {
		h.cfg.Logger().Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Secret updated successfully"})
}

func (h *SecretsHandlerImpl) DeleteSecret(ctx *gin.Context) {
	secretIDStr := ctx.Param("secretID")
	secretID, err := strconv.Atoi(secretIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.serviceProvider.DeleteSecret(ctx, secretID); err != nil {
		h.cfg.Logger().Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Secret deleted successfully"})
}
