package infrastructure

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Ozoniuss/hangcounts/config"
	"github.com/Ozoniuss/hangcounts/domain/model"
	"github.com/Ozoniuss/hangcounts/domain/storage"
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

func (p *PostgresStore) StoreIndividual(ctx context.Context, individual model.Individual, createdAt, updatedAt, deletedAt *time.Time) error {
	query := `
		INSERT INTO individuals (id, name, email, username, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	result, err := p.conn.Exec(ctx, query, individual.Id, individual.Name, individual.Email, individual.Username, createdAt, updatedAt, deletedAt)
	if err != nil {
		return fmt.Errorf("lol")
	}

	rows := result.RowsAffected()
	if rows != 1 {
		p.logger.ErrorContext(ctx, "expected one row to be affected when storing individual", slog.Int64("rows_affected", rows))
		return storage.ErrUnknown
	}

	return err
}
