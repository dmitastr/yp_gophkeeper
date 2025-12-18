package postgresdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"gophkeep/internal/core/models"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gophkeep/internal/config"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
	cfg  config.ConfigProvider
}

func NewPostgresStorage(ctx context.Context, cfg config.ConfigProvider) (*PostgresStorage, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.GetConfig().DBConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db with url=%s: %w", cfg.GetConfig().DBConnStr, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db with url=%s: %w", cfg.GetConfig().DBConnStr, err)
	}
	cfg.Logger().Info("Database connection established successfully")

	m, err := migrate.New(
		"file://database/migrations",
		cfg.GetConfig().DBConnStr)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	cfg.Logger().Info("Database migration succeeded")

	return &PostgresStorage{pool: pool, cfg: cfg}, nil
}

func (p *PostgresStorage) AddUser(ctx context.Context, user *models.User) (*models.User, error) {
	p.cfg.Logger().Info("RegisterUser called", zap.String("username", user.Username), zap.String("password", user.Password))
	query := `INSERT INTO users (username, password_hash, created_at) VALUES (@username, @password_hash, @created_at) RETURNING user_id`

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}

	var userID models.UserID
	if err := tx.QueryRow(ctx, query, user.ToNamedArgs()).Scan(&userID); err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("could not add user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}
	p.cfg.Logger().Info("Successfully add user")
	user.ID = userID

	return user, nil
}

func (p *PostgresStorage) GetUser(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	p.cfg.Logger().Info("GetUser called", zap.String("username", username))
	query := `SELECT user_id, username, password_hash, created_at FROM users WHERE username = $1`

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}

	if err := tx.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Hash, &user.CreatedAt); err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("could not add user: %w", err)
	}
	tx.Commit(ctx)

	p.cfg.Logger().Info("Successfully get user")

	return &user, nil
}

func (p *PostgresStorage) AddPassword(ctx context.Context, login string, password *models.Password) error {
	p.cfg.Logger().Info("AddPassword called", zap.String("login", login), zap.String("password", password.Password))

	return nil
}

func (p *PostgresStorage) GetPassword(ctx context.Context, login string, passwordID int) (*models.Password, error) {
	p.cfg.Logger().Info("AddPassword called", zap.String("login", login), zap.Int("passwordID", passwordID))

	return nil, nil
}

func (p *PostgresStorage) AddSecret(ctx context.Context, userID models.UserID, secretType models.SecretType, secret []byte, comment string) error {
	p.cfg.Logger().Info("AddSecret called", zap.String("secretType", string(secretType)), zap.Int("userID", int(userID)))
	query := `INSERT INTO secrets (user_id, secret_type, secret, created_at, comment) VALUES ($1, $2, $3, $4, $5)`

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	now := time.Now()
	if _, err := tx.Exec(ctx, query, userID, secretType, secret, now, comment); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("could not add secret: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (p *PostgresStorage) GetSecret(ctx context.Context, secretID int, userID models.UserID) (*models.Secret, error) {
	p.cfg.Logger().Info("GetSecret called", zap.Int("userID", int(userID)), zap.Int("secretID", secretID))
	query := `SELECT id, secret_type, created_at, secret, comment FROM secrets 
            WHERE user_id = $1 AND id = $2`

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}

	var secret models.Secret
	err = tx.QueryRow(ctx, query, userID, secretID).Scan(&secret.ID, &secret.Type, &secret.CreatedAt, &secret.Content, &secret.Comment)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("could not add user: %w", err)
	}
	tx.Commit(ctx)

	p.cfg.Logger().Info("Successfully get secrets")

	return &secret, nil
}

func (p *PostgresStorage) GetAllSecrets(ctx context.Context, userID models.UserID) ([]models.SecretInfo, error) {
	p.cfg.Logger().Info("GetAllSecrets called", zap.Int("userID", int(userID)))
	query := `SELECT id, secret_type, created_at FROM secrets WHERE user_id = $1
			ORDER BY created_at DESC`

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}

	rows, err := tx.Query(ctx, query, userID)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("could not add user: %w", err)
	}

	p.cfg.Logger().Info("Successfully get secrets")
	tx.Commit(ctx)

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.SecretInfo])
}
