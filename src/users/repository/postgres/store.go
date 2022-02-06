package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

// Stores a new user into the DB
func (r PostgresRepository) Store(ctx context.Context, user domain.User) (*domain.User, error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	query := `INSERT INTO users(id, first_name, last_name, username, password, role_id, created_at, updated_at) 
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, first_name, last_name, username, password, role_id, refresh_token_id, created_at, updated_at`

	statement, err := r.Db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	row := statement.
		QueryRowContext(ctx,
			user.ID,
			user.FirstName,
			user.LastName,
			user.Username,
			user.Password,
			user.RoleId,
			time.Now(),
			time.Now(),
		)

	return r.scanUserRow(row)
}
