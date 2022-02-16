package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/plagioriginal/user-microservice/database/migrations"
	_refreshTokensMigrations "github.com/plagioriginal/user-microservice/refresh-tokens/migrations"
	_rolesMigrations "github.com/plagioriginal/user-microservice/roles/migrations"
	_usersMigrations "github.com/plagioriginal/user-microservice/users/migrations"
)

type MigrationSettings struct {
	DefaultUserUsername string
	DefaultUserPassword string
	JwtSecret           string
	Timeout             time.Duration
	BcryptCost          int
}

func DoMigrations(l *log.Logger, db *sql.DB, settings MigrationSettings) {
	migrations := migrations.MigrationsHandler{
		Logger: l,
		Db:     db,
		Migrations: []migrations.Migration{
			// Migrations
			migrations.NewAddUuidExtensionMigration(),
			_rolesMigrations.NewCreateRolesMigration(),
			_usersMigrations.NewCreateUsersMigration(),
			_refreshTokensMigrations.NewCreateRefreshTokensMigration(),
			_usersMigrations.NewAddRefreshTokenReferenceMigration(),

			// Seeds
			_rolesMigrations.NewAddRolesMigration(),
			_usersMigrations.NewAddDefaultUserMigration(settings.BcryptCost),
		},
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), settings.Timeout)
	ctx = context.WithValue(ctx, _usersMigrations.DefaultUserNameKey, settings.DefaultUserUsername)
	ctx = context.WithValue(ctx, _usersMigrations.DefaultUserPasswordKey, settings.DefaultUserPassword)
	ctx = context.WithValue(ctx, "jwtSecret", settings.JwtSecret)

	defer cancelfunc()
	migrations.DoAll(ctx)
}
