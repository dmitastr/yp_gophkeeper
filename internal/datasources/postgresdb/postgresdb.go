package postgresdb

import (
	"context"
	"fmt"

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
	//
	// m, err := migrate.New(
	// 	"file://database/migrations",
	// 	cfg.DatabaseURI)
	// if err != nil {
	// 	return nil, err
	// }
	// if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
	// 	return nil, err
	// }
	// logger.Infof("Database migration succeeded")
	//
	// rp := retrypolicy.NewRetryPolicy(3, time.Second, true)

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db with url=%s: %w", cfg.GetConfig().DBConnStr, err)
	}
	cfg.Logger().Info("Database connection established successfully")

	return &PostgresStorage{pool: pool}, nil
}

func (p *PostgresStorage) AddUser(username string, password string) error {
	p.cfg.Logger().Info("AddUser called", zap.String("username", username), zap.String("password", password))

	return nil
}

func (p *PostgresStorage) GetUser(username string) (*models.User, error) {
	p.cfg.Logger().Info("GetUser called", zap.String("username", username))

	return nil, nil
}

func (p *PostgresStorage) AddPassword(login string, password *models.Password) error {
	p.cfg.Logger().Info("AddPassword called", zap.String("login", login), zap.String("password", password.Password))

	return nil
}

func (p *PostgresStorage) GetPassword(login string, passwordID int) (*models.Password, error) {
	p.cfg.Logger().Info("AddPassword called", zap.String("login", login), zap.Int("passwordID", passwordID))

	return nil, nil
}
