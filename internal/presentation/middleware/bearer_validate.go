package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/domain/service/authservice"
)

type BearerValidator interface {
	VerifyJWT(c *gin.Context)
}

type BearerValidatorImpl struct {
	cfg         config.ConfigProvider
	authService authservice.AuthService
}

func NewBearerValidator(cfg config.ConfigProvider, authService authservice.AuthService) *BearerValidatorImpl {
	return &BearerValidatorImpl{cfg: cfg, authService: authService}
}

func (b *BearerValidatorImpl) VerifyJWT(c *gin.Context) {
	bearerToken := c.Request.Header.Get("Authorization")
	reqToken := strings.Split(bearerToken, " ")[1]

	claims, err := b.authService.VerifyJWT(reqToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		b.cfg.Logger().Error("error while validating token", zap.Error(err))
		return
	}

	username, _ := claims.GetSubject()
	issuer, _ := claims.GetIssuer()
	userID := claims.UserID

	b.cfg.Logger().Info("Request from", zap.Int("userID", int(userID)), zap.String("username", username))

	c.Set("username", username)
	c.Set("issuer", issuer)
	c.Set("userID", userID)

	c.Next()
}
