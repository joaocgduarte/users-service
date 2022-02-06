package postgres

import (
	"context"

	"github.com/plagioriginal/user-microservice/domain"
)

// Gets a user by their respective username.
func (r PostgresRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, username, password, role_id, refresh_token_id, created_at, updated_at
		FROM users 
		WHERE username = $1 AND deleted_at IS NULL
		LIMIT 1
	`

	statement, err := r.Db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	row := statement.QueryRowContext(ctx, username)
	return r.scanUserRow(row)
}
