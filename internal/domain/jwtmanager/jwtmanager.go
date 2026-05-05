package jwtmanager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gophkeep/internal/config"
	"gophkeep/internal/core/models"
)

type Manager interface {
	IssueJWT(*models.User) (string, error)
	VerifyJWT(string) (*Claims, error)
}

type Claims struct {
	UserID models.UserID `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	cfg config.ConfigProvider
}

func New(cfg config.ConfigProvider) *JWTManager {
	return &JWTManager{cfg: cfg}
}

func (manager *JWTManager) IssueJWT(user *models.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gophkeeper",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(manager.cfg.GetConfig().Key))
	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}
	return tokenString, nil
}

func (manager *JWTManager) VerifyJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(manager.cfg.GetConfig().Key), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
