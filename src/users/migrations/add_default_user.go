package migrations

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/plagioriginal/user-microservice/database/migrations"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/helpers"
	_rolesRepo "github.com/plagioriginal/user-microservice/roles/repository/postgres"
	_usersRepo "github.com/plagioriginal/user-microservice/users/repository/postgres"
	_usersService "github.com/plagioriginal/user-microservice/users/service"
)

const (
	DefaultUserNameKey     string = "defaultUserUsername"
	DefaultUserPasswordKey string = "defaultUserPassword"
)

// Adds default user to the DB.
func AddDefaultUser(bcryptCost int) func(ctx context.Context, db *sql.DB, logger *log.Logger) error {
	return func(ctx context.Context, db *sql.DB, logger *log.Logger) error {
		defaultUserUsername := helpers.StringFromContext(ctx, DefaultUserNameKey)
		defaultUserPassword := helpers.StringFromContext(ctx, DefaultUserPasswordKey)

		if len(defaultUserUsername) == 0 || len(defaultUserPassword) == 0 {
			return errors.New("cant add default users without credentials")
		}

		timeoutDuration := time.Duration(5) * time.Second
		roleRepo := _rolesRepo.New(db)
		userRepo := _usersRepo.New(db)
		adminRoleSlug := domain.DEFAULT_ROLE_ADMIN.RoleSlug

		userService := _usersService.New(
			userRepo,
			roleRepo,
			timeoutDuration,
			bcryptCost,
		)

		user, _ := userRepo.GetByUsername(ctx, defaultUserUsername)

		if user != nil {
			logger.Println("user already exists, skipping add-default-user migration")
			return nil
		}

		req := domain.StoreUserRequest{
			Username: defaultUserUsername,
			Password: defaultUserPassword,
			RoleSlug: adminRoleSlug,
		}
		if _, err := userService.Store(ctx, req); err != nil {
			logger.Printf("error storing user on add-default-user migration: %v", err)
			return err
		}
		return nil
	}
}

// Adds default user to the DB.
func NewAddDefaultUserMigration(bcryptCost int) migrations.Migration {
	return migrations.Migration{
		Name: "add-default-user",
		Up:   AddDefaultUser(bcryptCost),
	}
}
