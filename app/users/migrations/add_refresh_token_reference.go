package migrations

import (
	"context"
	"database/sql"

	"github.com/plagioriginal/user-microservice/database/migrations"
)

// Creates the users table
func AddRefreshTokenReference(ctx context.Context, db *sql.DB) error {
	query := `
		ALTER TABLE users
		ADD COLUMN IF NOT EXISTS refresh_token_id uuid DEFAULT NULL REFERENCES refresh_tokens ON DELETE CASCADE ON UPDATE CASCADE
	`

	_, err := db.ExecContext(ctx, query)

	return err
}

// Creates a new migration
func NewAddRefreshTokenReferenceMigration() migrations.Migration {
	return migrations.Migration{
		Name: "add-refresh-token-reference",
		Up:   AddRefreshTokenReference,
	}
}
