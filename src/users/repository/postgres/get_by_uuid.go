package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

// Gets a user by uuid
func (r PostgresRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, username, password, role_id, refresh_token_id, created_at, updated_at
		FROM users 
		WHERE id = $1  AND deleted_at IS NULL
		LIMIT 1
	`

	statement, err := r.Db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	row := statement.QueryRowContext(ctx, uuid)
	return r.scanUserRow(row)
}
