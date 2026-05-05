package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gophkeep/internal/config"
	"gophkeep/internal/domain/service"
)

// PingHandler is an interface for Ping method
type PingHandler interface {
	Ping(*gin.Context)
}

// HandlerImpl implements [PingHandler]
type HandlerImpl struct {
	serviceProvider service.IService
	cfg             config.ConfigProvider
}

// NewHandler creates new [PingHandler]
func NewHandler(cfg config.ConfigProvider, serviceProvider service.IService) PingHandler {
	return &HandlerImpl{serviceProvider: serviceProvider, cfg: cfg}
}

// Ping checks is server is running and returns username and app name
func (h *HandlerImpl) Ping(c *gin.Context) {
	if err := h.serviceProvider.Ping(c); err != nil {
		h.cfg.Logger().Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")
	issuer := c.GetString("issuer")

	c.JSON(http.StatusOK, gin.H{"username": username, "issuer": issuer})
}
