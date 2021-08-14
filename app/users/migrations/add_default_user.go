package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/plagioriginal/user-microservice/database/migrations"
	"github.com/plagioriginal/user-microservice/domain"
	_rolesRepo "github.com/plagioriginal/user-microservice/roles/repository/postgres"
	_usersRepo "github.com/plagioriginal/user-microservice/users/repository/postgres"
	"github.com/plagioriginal/user-microservice/users/service"
	"github.com/plagioriginal/user-microservice/users/tokens"
)

// Adds default user to the DB.
func AddDefaultUser(ctx context.Context, db *sql.DB) error {
	defaultUserUsername, defaultUserPassword := ctx.Value("defaultUserUsername").(string), ctx.Value("defaultUserPassword").(string)

	if len(defaultUserUsername) == 0 || len(defaultUserPassword) == 0 {
		return errors.New("cant add default users without credentials")
	}

	roleRepo := _rolesRepo.New(db)
	userRepo := _usersRepo.New(db)
	adminRoleSlug := domain.DEFAULT_ROLE_ADMIN.RoleSlug

	userService := service.New(
		userRepo,
		roleRepo,
		tokens.NewTokenManager(ctx.Value("jwtSecret").(string)),
		time.Duration(2)*time.Second,
	)

	userService.Store(ctx, defaultUserUsername, defaultUserPassword, adminRoleSlug)
	return nil
}

// Adds default user to the DB.
func NewAddDefaultUserMigration() migrations.Migration {
	return migrations.Migration{
		Name: "add-default-user",
		Up:   AddDefaultUser,
	}
}
