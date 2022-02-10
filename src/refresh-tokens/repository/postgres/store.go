package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

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
