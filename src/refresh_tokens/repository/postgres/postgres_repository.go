package postgres

import (
	"database/sql"

	"github.com/plagioriginal/user-microservice/domain"
)

type PostgresRepository struct {
	Db *sql.DB
}

func New(db *sql.DB) domain.RefreshTokenRepository {
	return PostgresRepository{db}
}
