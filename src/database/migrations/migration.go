package migrations

import (
	"context"
	"database/sql"
	"log"
)

// The single migration
type Migration struct {
	Name string
	Up   func(ctx context.Context, db *sql.DB, logger *log.Logger) error
}
