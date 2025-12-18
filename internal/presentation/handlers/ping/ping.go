package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gophkeep/internal/config"
	"gophkeep/internal/domain/service"
)

type PingHandler interface {
	Ping(*gin.Context)
}

type HandlerImpl struct {
	serviceProvider service.IService
	cfg             config.ConfigProvider
}

func NewHandler(cfg config.ConfigProvider, serviceProvider service.IService) PingHandler {
	return &HandlerImpl{serviceProvider: serviceProvider, cfg: cfg}
}

func (h *HandlerImpl) Ping(c *gin.Context) {
	username := c.GetString("username")
	issuer := c.GetString("issuer")

	c.JSON(http.StatusOK, gin.H{"username": username, "issuer": issuer})
}
