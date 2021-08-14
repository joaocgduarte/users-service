package migrations

import (
	"context"
	"database/sql"

	"github.com/plagioriginal/user-microservice/database/migrations"
)

// Creates the refesh tokens table
func CreateRefreshTokensTable(ctx context.Context, db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS refresh_tokens(
			id uuid DEFAULT uuid_generate_v4() NOT NULL,
			token uuid DEFAULT uuid_generate_v4() NOT NULL,
			valid_until timestamptz NOT NULL DEFAULT (NOW() + INTERVAL '7 days'),
			PRIMARY KEY (id)
		);
	`

	_, err := db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, query)

	return err
}

// Creates a new migration
func NewCreateRefreshTokensMigration() migrations.Migration {
	return migrations.Migration{
		Name: "create-roles",
		Up:   CreateRefreshTokensTable,
	}
}
