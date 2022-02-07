package migrations

import (
	"context"
	"database/sql"
	"log"

	"github.com/plagioriginal/user-microservice/database/migrations"
	"github.com/plagioriginal/user-microservice/domain"
	roleRepo "github.com/plagioriginal/user-microservice/roles/repository/postgres"
)

// Adds some default roles.
func AddDefaultRoles(ctx context.Context, db *sql.DB, logger *log.Logger) error {

	repo := roleRepo.New(db)

	roles := []domain.Role{
		domain.DEFAULT_ROLE_ADMIN,
		domain.DEFAULT_ROLE_USER,
	}

	for _, role := range roles {
		logger.Printf("adding role %v\n", role)
		repo.Store(ctx, &role)
	}

	return nil
}

// Creates a new migration
func NewAddRolesMigration() migrations.Migration {
	return migrations.Migration{
		Name: "add-default-roles",
		Up:   AddDefaultRoles,
	}
}
