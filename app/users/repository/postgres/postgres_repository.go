package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

// Definition of the repository
type PostgresRepository struct {
	Db *sql.DB
}

// Instantiates the repository
func New(db *sql.DB) domain.UserRepository {
	return PostgresRepository{db}
}

// Scans a user row without the deleted_at value
func (r PostgresRepository) scanUserRow(row *sql.Row) (*domain.User, error) {
	result := domain.User{}

	err := row.Scan(
		&result.ID,
		&result.FirstName,
		&result.LastName,
		&result.Username,
		&result.Password,
		&result.RoleId,
		&result.RefreshTokenId,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return &domain.User{}, err
	}

	return &result, nil
}

// Lists the users in the DB in a paginated matter.
func (r PostgresRepository) List(ctx context.Context, page int, perPage int) ([]domain.User, error) {
	panic("to be implemented")
}

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
		return &domain.User{}, err
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
		return &domain.User{}, err
	}

	row := statement.QueryRowContext(ctx, uuid)
	return r.scanUserRow(row)
}

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
		return &domain.User{}, err
	}

	row := statement.QueryRowContext(ctx, username)
	return r.scanUserRow(row)
}
