package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Ozoniuss/hangcounts/config"
	"github.com/Ozoniuss/hangcounts/domain/model"
	"github.com/Ozoniuss/hangcounts/domain/storage"
	"github.com/jackc/pgx/v5"
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

func (p *PostgresStore) StoreIndividual(ctx context.Context, individual model.Individual) error {
	query := `
		INSERT INTO individuals (name, email, username, created_at)
		VALUES ($1, $2, $3, $4);
	`

	createdAt := time.Now()
	result, err := p.conn.Exec(ctx, query, individual.Name, individual.Email, individual.Username, createdAt)
	if err != nil {
		p.logger.ErrorContext(ctx, "failed to execute query", slog.Any("error", err))
		return storage.ErrUnknown
	}

	rows := result.RowsAffected()
	if rows != 1 {
		p.logger.ErrorContext(ctx, "expected one row to be affected", slog.Int64("rows_affected", rows))
		return storage.ErrUnknown
	}

	return nil
}

func (p *PostgresStore) GetIndividual(ctx context.Context, individualUsername model.IndividualId) (model.Individual, error) {
	query := `
		SELECT username, email, name, created_at, updated_at, deleted_at
		FROM individuals
		WHERE username=$1;
	`

	row := p.conn.QueryRow(ctx, query, individualUsername)
	var username string
	var name string
	var email string
	var created_at, updated_at time.Time

	var deleted_at sql.NullTime

	err := row.Scan(&username, &email, &name, &created_at, &updated_at, &deleted_at)
	if errors.Is(err, pgx.ErrNoRows) {
		p.logger.ErrorContext(ctx, "individual not found", slog.String("username", string(individualUsername)))
		return model.Individual{}, storage.ErrNotFound
	} else if err != nil {
		p.logger.ErrorContext(ctx, "unknown error when retrieving individual", slog.String("username", string(individualUsername)), slog.Any("error", err))
		return model.Individual{}, storage.ErrUnknown
	}

	if deleted_at.Valid {
		return model.Individual{}, storage.ErrDeleted
	}

	return model.Individual{
		Username: model.IndividualId(username),
		Name:     name,
		Email:    model.Email(email),
	}, nil
}

func (p *PostgresStore) DeleteIndividual(ctx context.Context, individualUsername model.IndividualId) error {

	selectQuery := `
		SELECT username, deleted_at
		FROM individuals
		WHERE username=$1;
	`

	row := p.conn.QueryRow(ctx, selectQuery, individualUsername)
	var username string
	var deleted_at sql.NullTime

	err := row.Scan(&username, &deleted_at)
	if errors.Is(err, pgx.ErrNoRows) {
		p.logger.ErrorContext(ctx, "individual not found", slog.String("username", string(individualUsername)))
		return storage.ErrNotFound
	} else if err != nil {
		p.logger.ErrorContext(ctx, "unknown error", slog.String("username", string(individualUsername)), slog.Any("error", err))
		return storage.ErrUnknown
	}

	// once deleted_at was populated do not allow any changes to the record
	if deleted_at.Valid {
		return storage.ErrDeleted
	}

	deleteQuery := `
		UPDATE individuals
		SET deleted_at=$2
		WHERE username=$1;
	`

	deletedAt := time.Now()
	result, err := p.conn.Exec(ctx, deleteQuery, individualUsername, deletedAt)
	if err != nil {
		p.logger.ErrorContext(ctx, "failed to execute query", slog.Any("error", err))
		return storage.ErrUnknown
	}

	rows := result.RowsAffected()
	if rows != 1 {
		p.logger.ErrorContext(ctx, "expected one row to be affected", slog.Int64("rows_affected", rows))
		return storage.ErrUnknown
	}

	return nil
}
