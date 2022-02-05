package migrations

import (
	"context"
	"database/sql"
	"log"

	"github.com/plagioriginal/user-microservice/database/migrations"
)

// Creates the roles table
func CreateRolesTable(ctx context.Context, db *sql.DB, logger *log.Logger) error {
	query := `CREATE TABLE IF NOT EXISTS roles(
		id uuid DEFAULT uuid_generate_v4(),
		role_slug varchar(255) NOT NULL UNIQUE,
		role_label varchar(255) NOT NULL UNIQUE,

		created_at timestamptz NOT NULL DEFAULT (now()),
		updated_at timestamptz NOT NULL DEFAULT (now()),
		deleted_at timestamptz DEFAULT NULL,
		PRIMARY KEY (id)
	);`

	_, err := db.ExecContext(ctx, query)

	return err
}

// Creates a new migration
func NewCreateRolesMigration() migrations.Migration {
	return migrations.Migration{
		Name: "create-roles-table",
		Up:   CreateRolesTable,
	}
}
