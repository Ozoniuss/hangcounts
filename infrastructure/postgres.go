package infrastructure

import (
	"context"
	"fmt"

	"github.com/Ozoniuss/hangcounts/config"
	"github.com/jackc/pgx/v5"
)

type PostgresStore struct {
	Conn *pgx.Conn
}

func NewPostgresStore(ctx context.Context, cfg config.PostgresConfig) (*PostgresStore, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName)
	connCfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("could not parse dsn: %w", err)
	}
	conn, err := pgx.ConnectConfig(ctx, connCfg)
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
