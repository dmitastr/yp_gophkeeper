package datasources

import (
	"context"

	"gophkeep/internal/core/models"
)

type Datasource interface {
	AddUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
	AddPassword(ctx context.Context, login string, password *models.Password) error
	GetPassword(ctx context.Context, login string, passwordID int) (*models.Password, error)
	AddSecret(ctx context.Context, userID models.UserID, secretType models.SecretType, secret []byte, comment string) error
	GetSecret(ctx context.Context, secretID int, userID models.UserID) (*models.Secret, error)
	UpdateSecret(ctx context.Context, secret *models.Secret, userID models.UserID) error
	DeleteSecret(ctx context.Context, secretID int, userID models.UserID) error
	GetAllSecrets(ctx context.Context, userID models.UserID) ([]models.SecretInfo, error)
	Ping(ctx context.Context) error
	Close() error
}
