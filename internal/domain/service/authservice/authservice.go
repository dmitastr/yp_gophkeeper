package authservice

import (
	"errors"
	"fmt"

	"gophkeep/internal/config"
	"gophkeep/internal/datasources"
	"gophkeep/internal/domain/jwtmanager"
	"gophkeep/internal/domain/models"
	"gophkeep/internal/presentation/params"
)

type AuthService interface {
	Authenticate(object params.AuthRequestObject) (string, error)
	VerifyJWT(string) (*jwtmanager.Claims, error)
}

type AuthServiceImpl struct {
	manager jwtmanager.Manager
	db      datasources.Datasource
}

func NewAuthService(cfg config.ConfigProvider, db datasources.Datasource) AuthService {
	manager := jwtmanager.New(cfg)
	return &AuthServiceImpl{manager: manager, db: db}
}

func (a *AuthServiceImpl) Authenticate(object params.AuthRequestObject) (string, error) {
	if object.Password == "" || object.Username == "" {
		return "", errors.New("invalid object")
	}
	user := &models.User{Username: object.Username}
	token, err := a.manager.IssueJWT(user)
	if err != nil {
		return "", fmt.Errorf("issue token: %w", err)
	}
	return token, nil
}

func (a *AuthServiceImpl) VerifyJWT(token string) (*jwtmanager.Claims, error) {
	return a.manager.VerifyJWT(token)
}
