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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CONSTRAINT_UNIQUE_INDIVIDUAL_USERNAME  = "unique_individual_usernamename"
	CONSTRAINT_UNIQUE_INDIVIDUAL_EMAIL     = "unique_individual_email"
	CONSTRAINT_UNIQUE_HANGOUT_PUBLIC_ID    = "unique_hangout_public_id"
	CONSTRAINT_FOREIGN_KEY_HANGOUT_CREATOR = "fk_hangout_creator"
	CONSTRAINT_FOREIGN_KEY_HANGOUT         = "fk_hangout"
	CONSTRAINT_FOREIGN_KEY_INDIVIDUAL      = "fk_individual"
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
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			p.logger.DebugContext(ctx, "got pgerr", slog.String("error", fmt.Sprintf("%#v", pgerr)))
			if pgerr.Code != pgerrcode.UniqueViolation {
				return storage.ErrUnknown
			}
			switch pgerr.ConstraintName {
			case CONSTRAINT_UNIQUE_INDIVIDUAL_EMAIL:
				return storage.ErrIndividualEmailAlreadyExists
			case CONSTRAINT_UNIQUE_INDIVIDUAL_USERNAME:
				return storage.ErrIndividualUsernameAlreadyExists
			}
		}
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

func (p *PostgresStore) MarkIndividualAsDeleted(ctx context.Context, individualUsername model.IndividualId) error {

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

func (p *PostgresStore) StoreHangoutOfIndividuals(ctx context.Context, hangout model.Hangout) error {

	queryHangoutIndividual := `
		SELECT id, deleted_at
		FROM individuals
		WHERE username = $1;
	`

	queryHangoutDetails := `
		INSERT INTO hangouts (public_id, location, description, duration_minutes, date, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	queryInsertParticipant := `
		INSERT INTO hangout_individuals (hangout_id, individual_id, created_at)
		VALUES ($1, $2, $3);
	`
	currentTimestamp := time.Now()

	tx, err := p.conn.Begin(ctx)
	if err != nil {
		p.logger.ErrorContext(ctx, "failed to start a transaction", slog.Any("error", err))
		return storage.ErrUnknown
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			p.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("error", err))
		}
	}()

	row := tx.QueryRow(ctx, queryHangoutIndividual, hangout.CreatedBy)
	var creatorId int
	var creatorDeleted sql.NullTime

	if err := row.Scan(&creatorId, &creatorDeleted); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.logger.ErrorContext(ctx, "hangout creator not found", slog.String("username", string(hangout.CreatedBy)))
			return storage.ErrHangoutCreatorNotFound
		}
		p.logger.ErrorContext(ctx, "unknown error when retrieving creator", slog.String("username", string(hangout.CreatedBy)), slog.Any("error", err))
		return storage.ErrUnknown
	}
	if creatorDeleted.Valid {
		return storage.ErrHangoutCreatorDeleted
	}

	p.logger.Debug("retrieved creator", slog.Int("id", int(creatorId)), slog.Any("deleted_at", creatorDeleted))

	var hangoutId int64
	row = tx.QueryRow(ctx, queryHangoutDetails, hangout.PublicId, hangout.Location, hangout.Description, hangout.Duration, hangout.Date, creatorId, currentTimestamp)
	if err = row.Scan(&hangoutId); err != nil {
		p.logger.ErrorContext(ctx, "failed to execute hangout creation query", slog.Any("error", err))
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			p.logger.DebugContext(ctx, "got pgerr", slog.String("error", fmt.Sprintf("%#v", pgerr)))
			if pgerr.Code != pgerrcode.UniqueViolation && pgerr.Code != pgerrcode.ForeignKeyViolation {
				return storage.ErrUnknown
			}
			switch pgerr.ConstraintName {
			case CONSTRAINT_UNIQUE_HANGOUT_PUBLIC_ID:
				return storage.ErrAlreadyExists
			// depending on isolation level, this may be redundant, but
			// in theory should never happen
			case CONSTRAINT_FOREIGN_KEY_HANGOUT_CREATOR:
				p.logger.WarnContext(ctx, "foreign key constraint violation in hangouts table")
				return storage.ErrHangoutCreatorNotFound
			}
			return storage.ErrUnknown
		}
	}
	p.logger.Debug("retrieved hangout", slog.Int64("hid", hangoutId))

	// the database does not enforce the hangout to have at least one
	// participant, nor does it enforce that the creator is part of the
	// participant.
	for _, participant := range hangout.Individuals {
		row := tx.QueryRow(ctx, queryHangoutIndividual, participant)

		var participantId int
		var participantDeleted sql.NullTime
		if err := row.Scan(&participantId, &participantDeleted); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				p.logger.ErrorContext(ctx, "hangout participant not found", slog.String("username", string(hangout.CreatedBy)))
				return storage.ErrHangoutParticipantNotFound
			}
			p.logger.ErrorContext(ctx, "unknown error when retrieving creator", slog.String("username", string(hangout.CreatedBy)), slog.Any("error", err))
			return storage.ErrUnknown
		}
		if participantDeleted.Valid {
			return storage.ErrHangoutParticipantDeleted
		}

		p.logger.Debug("retrieved participant", slog.Int("id", int(participantId)), slog.Any("deleted_at", participantDeleted))

		result, err := tx.Exec(ctx, queryInsertParticipant, hangoutId, participantId, currentTimestamp)
		// depending on the isolation level, it might not be necessary to check those
		if err != nil {
			p.logger.ErrorContext(ctx, "failed to execute participant insertion query", slog.Any("error", err))
			var pgerr *pgconn.PgError
			if errors.As(err, &pgerr) {
				p.logger.DebugContext(ctx, "got pgerr", slog.String("error", fmt.Sprintf("%#v", pgerr)))
				if pgerr.Code != pgerrcode.ForeignKeyViolation {
					return storage.ErrUnknown
				}
				// these may be redundant depending on isolation level, but
				// in theory should never happen
				switch pgerr.ConstraintName {
				case CONSTRAINT_FOREIGN_KEY_HANGOUT:
					p.logger.WarnContext(ctx, "foreign key constraint violation in hangout_individuals table", slog.String("constraint", CONSTRAINT_FOREIGN_KEY_HANGOUT))
					return storage.ErrParticipantHangoutNotFound
				case CONSTRAINT_FOREIGN_KEY_INDIVIDUAL:
					p.logger.WarnContext(ctx, "foreign key constraint violation in hangout_individuals table", slog.String("constraint", CONSTRAINT_FOREIGN_KEY_INDIVIDUAL))
					return storage.ErrParticipantIndividualNotFound
				}
			}
			return storage.ErrUnknown
		}

		rows := result.RowsAffected()
		if rows != 1 {
			p.logger.ErrorContext(ctx, "expected one row to be affected", slog.Int64("rows_affected", rows))
			return storage.ErrUnknown
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.logger.ErrorContext(ctx, "could not commit transaction", slog.Any("error", err))
	}

	p.logger.InfoContext(ctx, "hangout created", slog.String("creator", string(hangout.CreatedBy)), slog.Any("participants", hangout.Individuals))

	return nil
}
