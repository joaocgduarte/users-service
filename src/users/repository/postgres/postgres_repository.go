package postgres

import (
	"context"
	"database/sql"

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
		return nil, err
	}

	return &result, nil
}

// Lists the users in the DB in a paginated matter.
func (r PostgresRepository) List(ctx context.Context, page int, perPage int) ([]domain.User, error) {
	panic("to be implemented")
}
