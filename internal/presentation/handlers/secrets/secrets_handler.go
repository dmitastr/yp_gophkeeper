package secrets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gophkeep/internal/config"
	"gophkeep/internal/domain/models"
	"gophkeep/internal/domain/service"
)

type SecretsHandler interface {
	AddPassword(ctx *gin.Context)
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
	}

	if err := h.serviceProvider.AddPassword(&request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password added successfully"})
}
