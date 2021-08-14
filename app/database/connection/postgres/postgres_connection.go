package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Settings for the connection
type PostgresConnectionSettings struct {
	User     string
	Password string
	DBName   string
	Host     string
	SSLMode  string
}

// Gets the DB connection
func Get(settings PostgresConnectionSettings) (*sql.DB, error) {
	sslMode := "disable"

	if len(settings.SSLMode) > 0 {
		sslMode = settings.SSLMode
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		settings.User,
		settings.Password,
		settings.Host,
		settings.DBName,
		sslMode,
	)

	return sql.Open("postgres", connStr)
}
