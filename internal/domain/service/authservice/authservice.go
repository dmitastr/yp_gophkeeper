package authservice

import (
	"context"
	"errors"
	"fmt"

	"gophkeep/internal/config"
	"gophkeep/internal/core/models"
	"gophkeep/internal/datasources"
	"gophkeep/internal/domain/hashvalidator"
	"gophkeep/internal/domain/jwtmanager"
	"gophkeep/internal/presentation/params"
)

type AuthService interface {
	LoginUser(object params.AuthRequestObject) (string, error)
	VerifyJWT(string) (*jwtmanager.Claims, error)
	RegisterUser(object params.AuthRequestObject) (string, error)
}

type AuthServiceImpl struct {
	manager jwtmanager.Manager
	db      datasources.Datasource
	hash    hashvalidator.HashValidator
}

func NewAuthService(cfg config.ConfigProvider, db datasources.Datasource) AuthService {
	manager := jwtmanager.New(cfg)
	return &AuthServiceImpl{manager: manager, db: db, hash: hashvalidator.NewHashValidator()}
}

func (a *AuthServiceImpl) LoginUser(object params.AuthRequestObject) (string, error) {
	if object.Password == "" || object.Username == "" {
		return "", errors.New("invalid object")
	}

	user := &models.User{Username: object.Username, Password: object.Password}
	userExisted, err := a.db.GetUser(context.TODO(), object.Username)
	if err != nil {
		return "", err
	}
	if userExisted == nil {
		return "", errors.New("user not found")
	}

	if ok := a.hash.Validate(object.Password, userExisted.Hash); !ok {
		return "", errors.New("invalid password")
	}

	token, err := a.manager.IssueJWT(user)
	if err != nil {
		return "", fmt.Errorf("issue token: %w", err)
	}
	return token, nil
}

func (a *AuthServiceImpl) RegisterUser(object params.AuthRequestObject) (string, error) {
	user := &models.User{Username: object.Username, Hash: a.hash.CalculateHash(object.Password)}
	userAdded, err := a.db.AddUser(context.TODO(), user)
	if err != nil {
		return "", fmt.Errorf("add user error: %w", err)
	}

	token, err := a.manager.IssueJWT(userAdded)
	if err != nil {
		return "", fmt.Errorf("issue token: %w", err)
	}
	return token, nil

}

func (a *AuthServiceImpl) VerifyJWT(token string) (*jwtmanager.Claims, error) {
	return a.manager.VerifyJWT(token)
}
