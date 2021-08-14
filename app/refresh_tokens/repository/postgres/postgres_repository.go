package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

type PostgresRepository struct {
	Db *sql.DB
}

func New(db *sql.DB) domain.RefreshTokenRepository {
	return PostgresRepository{db}
}

// Gets a token
func (r PostgresRepository) GetByToken(ctx context.Context, token uuid.UUID) (domain.RefreshToken, error) {
	query := `
		SELECT id, token, valid_until
		FROM refresh_tokens
		WHERE token = $1
	`

	stmt, err := r.Db.PrepareContext(ctx, query)

	if err != nil {
		return domain.RefreshToken{}, err
	}

	row := stmt.QueryRowContext(ctx, token, time.Now())

	result := domain.RefreshToken{}

	err = row.Scan(
		&result.Id,
		&result.Token,
		&result.ValidUntil,
	)

	if err != nil {
		return domain.RefreshToken{}, err
	}

	return result, nil
}

// Stores the token into the DB
func (r PostgresRepository) Store(ctx context.Context, token domain.RefreshToken) (domain.RefreshToken, error) {
	query := `
		INSERT INTO refresh_tokens(id, token, valid_until)
		VALUES ($1, $2, $3)
		RETURNING id, token, valid_until
	`

	stmt, err := r.Db.PrepareContext(ctx, query)

	if err != nil {
		return domain.RefreshToken{}, err
	}

	if token.Id == uuid.Nil {
		token.Id = uuid.New()
	}

	result := domain.RefreshToken{}
	row := stmt.QueryRowContext(ctx, token.Id, token.Token, token.ValidUntil)
	err = row.Scan(
		&result.Id,
		&result.Token,
		&result.ValidUntil,
	)

	if err != nil {
		return domain.RefreshToken{}, err
	}

	return result, nil
}

func (r PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE id = $1
	`

	stmt, err := r.Db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, id)

	return err
}
