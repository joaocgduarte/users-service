package migrations

import (
	"context"
	"database/sql"

	"github.com/plagioriginal/user-microservice/database/migrations"
)

// Creates the users table
func CreateUsersTable(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id uuid DEFAULT uuid_generate_v4(),
		first_name varchar(255) DEFAULT '',
		last_name varchar(255) DEFAULT '',
		username varchar(255) NOT NULL UNIQUE,
		password text NOT NULL,
		role_id uuid NOT NULL REFERENCES roles ON DELETE CASCADE ON UPDATE CASCADE,

		created_at timestamptz NOT NULL DEFAULT (now()),
		updated_at timestamptz NOT NULL DEFAULT (now()),
		deleted_at timestamptz DEFAULT NULL,

		PRIMARY KEY (id),
		CONSTRAINT fk_users
		FOREIGN KEY(role_id) 
			REFERENCES roles(id)
	);`

	_, err := db.ExecContext(ctx, query)

	return err
}

// Creates a new migration
func NewCreateUsersMigration() migrations.Migration {
	return migrations.Migration{
		Name: "create-users-table",
		Up:   CreateUsersTable,
	}
}
