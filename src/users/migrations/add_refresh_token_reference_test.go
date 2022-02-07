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

func TestAddRefreshTokenReference_FailExec(t *testing.T) {
	migration := NewAddRefreshTokenReferenceMigration()
	assert.Equal(t, migration.Name, "add-refresh-token-reference")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `
		ALTER TABLE IF EXISTS users
		ADD COLUMN IF NOT EXISTS refresh_token_id uuid DEFAULT NULL REFERENCES refresh_tokens(id) ON DELETE SET NULL ON UPDATE CASCADE
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnError(errors.New("boom"))

	err = migration.Up(context.TODO(), db, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func TestAddRefreshTokenReference_TimeoutReached(t *testing.T) {
	migration := NewAddRefreshTokenReferenceMigration()
	assert.Equal(t, migration.Name, "add-refresh-token-reference")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `
		ALTER TABLE IF EXISTS users
		ADD COLUMN IF NOT EXISTS refresh_token_id uuid DEFAULT NULL REFERENCES refresh_tokens(id) ON DELETE SET NULL ON UPDATE CASCADE
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

func TestAddRefreshTokenReference_Success(t *testing.T) {
	migration := NewAddRefreshTokenReferenceMigration()
	assert.Equal(t, migration.Name, "add-refresh-token-reference")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `
		ALTER TABLE IF EXISTS users
		ADD COLUMN IF NOT EXISTS refresh_token_id uuid DEFAULT NULL REFERENCES refresh_tokens(id) ON DELETE SET NULL ON UPDATE CASCADE
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
