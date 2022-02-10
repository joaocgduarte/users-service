package migrations

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateRefreshTokens_FailExec(t *testing.T) {
	migration := NewCreateRefreshTokensMigration()
	assert.Equal(t, migration.Name, "create-refresh-tokens-table")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `
		CREATE TABLE IF NOT EXISTS refresh_tokens(
			id uuid DEFAULT uuid_generate_v4() NOT NULL,
			token uuid DEFAULT uuid_generate_v4() NOT NULL,
			valid_until timestamptz NOT NULL DEFAULT (NOW() + INTERVAL '7 days'),
			PRIMARY KEY (id)
		);
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnError(errors.New("boom"))

	err = migration.Up(context.TODO(), db, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func TestCreateRefreshTokens_TimeoutReached(t *testing.T) {
	migration := NewCreateRefreshTokensMigration()
	assert.Equal(t, migration.Name, "create-refresh-tokens-table")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `
		CREATE TABLE IF NOT EXISTS refresh_tokens(
			id uuid DEFAULT uuid_generate_v4() NOT NULL,
			token uuid DEFAULT uuid_generate_v4() NOT NULL,
			valid_until timestamptz NOT NULL DEFAULT (NOW() + INTERVAL '7 days'),
			PRIMARY KEY (id)
		);
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WillDelayFor(time.Duration(200 * time.Millisecond)).
		WillReturnError(errors.New("boom"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Millisecond))
	defer cancel()

	err = migration.Up(ctx, db, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
}

func TestCreateRefreshTokens_Success(t *testing.T) {
	migration := NewCreateRefreshTokensMigration()
	assert.Equal(t, migration.Name, "create-refresh-tokens-table")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `
		CREATE TABLE IF NOT EXISTS refresh_tokens(
			id uuid DEFAULT uuid_generate_v4() NOT NULL,
			token uuid DEFAULT uuid_generate_v4() NOT NULL,
			valid_until timestamptz NOT NULL DEFAULT (NOW() + INTERVAL '7 days'),
			PRIMARY KEY (id)
		);
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WillDelayFor(time.Duration(50 * time.Millisecond)).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Millisecond))
	defer cancel()

	err = migration.Up(ctx, db, nil)
	assert.Nil(t, err)
}
