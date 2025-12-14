package postgresdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gophkeep/internal/config"
	"gophkeep/internal/domain/models"
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

	return &PostgresStorage{pool: pool}, nil
}

func (p *PostgresStorage) AddUser(ctx context.Context, user *models.User) error {
	p.cfg.Logger().Info("AddUser called", zap.String("username", user.Username), zap.String("password", user.Password))
	query := `INSERT INTO users (username, password_hash, created_at) VALUES (@name, @hash, @created_at)`

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	if _, err := tx.Exec(ctx, query, user.ToNamedArgs()); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("could not add user: %w", err)
	}

	p.cfg.Logger().Info("Successfully add user")

	return nil
}

func (p *PostgresStorage) GetUser(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	p.cfg.Logger().Info("GetUser called", zap.String("username", username))
	query := `SELECT INTO user_id, username, password_hash, created_at FROM users WHERE username = $1`

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}

	if err := tx.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Hash, &user.CreatedAt); err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("could not add user: %w", err)
	}

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
