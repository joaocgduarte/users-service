package migrations

import (
	"context"
	"database/sql"
	"log"
)

// Wrapper for all the migrations.
type MigrationsHandler struct {
	Logger     *log.Logger
	Db         *sql.DB
	Migrations []Migration
}

// Does all the migrations
func (mh MigrationsHandler) DoAll(ctx context.Context) {
	for _, migration := range mh.Migrations {
		mh.Logger.Println("Migrating `" + migration.Name + "`...")
		err := migration.Up(ctx, mh.Db, mh.Logger)

		if err != nil {
			mh.Logger.Fatal(err)
		}

		mh.Logger.Println("Success.")
	}
}
