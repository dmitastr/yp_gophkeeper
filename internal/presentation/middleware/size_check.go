package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gophkeep/internal/config"
)

type SizeChecker interface {
	SizeCheck(c *gin.Context)
}

type sizeChecker struct {
	cfg config.ConfigProvider
}

func NewSizeChecker(cfg config.ConfigProvider) SizeChecker {
	return &sizeChecker{cfg: cfg}
}

func (s *sizeChecker) SizeCheck(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(
		c.Writer,
		c.Request.Body,
		s.cfg.GetConfig().MaxSize,
	)
	c.Next()
}
