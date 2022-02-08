package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

// Gets a single token by UUID
func (r PostgresRepository) GetByUUID(ctx context.Context, id uuid.UUID) (domain.RefreshToken, error) {
	query := `
		SELECT id, token, valid_until
		FROM refresh_tokens
		WHERE id = $1
	`

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return domain.RefreshToken{}, err
	}

	row := stmt.QueryRowContext(ctx, id)
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
