package migrations

import (
	"context"
	"database/sql"
)

// The single migration
type Migration struct {
	Name string
	Up   func(ctx context.Context, db *sql.DB) error
}
