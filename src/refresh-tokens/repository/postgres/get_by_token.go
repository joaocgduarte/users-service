package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

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

	row := stmt.QueryRowContext(ctx, token)
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
