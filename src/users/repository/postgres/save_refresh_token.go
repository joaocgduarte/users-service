package postgres

import (
	"context"

	"github.com/plagioriginal/user-microservice/domain"
)

// Saves a refresh token into the columns that has that link
func (r PostgresRepository) SaveRefreshToken(ctx context.Context, user *domain.User, token domain.RefreshToken) error {
	query := `
		UPDATE users
		SET refresh_token_id = $1
		WHERE id = $2
	`

	stmt, err := r.Db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, token.Id, user.ID)

	if err != nil {
		return err
	}

	return nil
}
