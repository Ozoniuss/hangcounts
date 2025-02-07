package infrastructure

import (
	"context"
	"fmt"

	"github.com/Ozoniuss/hangcounts/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	Conn *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, cfg config.PostgresConfig) (*PostgresStore, error) {

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

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStore{
		Conn: conn,
	}, nil
}
