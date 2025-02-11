package infrastructure

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Ozoniuss/hangcounts/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	conn   *pgxpool.Pool
	logger *slog.Logger
}

func NewPostgresStore(ctx context.Context, cfg config.PostgresConfig, logger *slog.Logger) (*PostgresStore, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName)
	connCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("could not parse dsn: %w", err)
	}
	conn, err := pgxpool.NewWithConfig(ctx, connCfg)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	logger.Debug("pool configuration", slog.Any("config", fmt.Sprintf("%+v", conn.Config())))

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStore{
		conn:   conn,
		logger: logger,
	}, nil
}

func (p *PostgresStore) StoreIndividual() {

}
