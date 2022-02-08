package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

type Repository struct {
	Db *sql.DB
}

func New(db *sql.DB) domain.RoleRepository {
	return Repository{db}
}

// Gets all the roles
func (r Repository) Fetch(ctx context.Context) ([]domain.Role, error) {
	result := make([]domain.Role, 0)

	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL;`
	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		role := domain.Role{}
		err = rows.Scan(
			&role.ID,
			&role.RoleSlug,
			&role.RoleLabel,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return make([]domain.Role, 0), err
		}

		result = append(result, role)
	}
	return result, nil
}

// Gets role by slug
func (r Repository) GetBySlug(ctx context.Context, slug string) (domain.Role, error) {
	result := domain.Role{}
	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL and role_slug=$1;`

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return result, err
	}

	row := stmt.QueryRowContext(ctx, slug)
	err = row.Scan(
		&result.ID,
		&result.RoleSlug,
		&result.RoleLabel,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	return result, err
}

// Gets role by UUID
func (r Repository) GetByUUID(ctx context.Context, uuid uuid.UUID) (domain.Role, error) {
	result := domain.Role{}
	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL and id=$1;`

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return result, err
	}

	row := stmt.QueryRowContext(ctx, uuid)
	err = row.Scan(
		&result.ID,
		&result.RoleSlug,
		&result.RoleLabel,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	return result, err
}

// Stores a new role into the DB
func (r Repository) Store(ctx context.Context, role domain.Role) (domain.Role, error) {
	result := domain.Role{}

	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}

	query := `
		INSERT INTO roles (id, role_slug, role_label, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, role_slug, role_label, created_at, updated_at
	`

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return result, err
	}

	row := stmt.QueryRowContext(ctx, role.ID, role.RoleSlug, role.RoleLabel, time.Now(), time.Now())

	err = row.Scan(
		&result.ID,
		&result.RoleSlug,
		&result.RoleLabel,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return domain.Role{}, err
	}

	return result, nil
}
