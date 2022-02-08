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
	migration := NewAddUuidExtensionMigration()
	assert.Equal(t, migration.Name, "add-uuid-extension")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnError(errors.New("boom"))

	err = migration.Up(context.TODO(), db, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func TestCreateRefreshTokens_TimeoutReached(t *testing.T) {
	migration := NewAddUuidExtensionMigration()
	assert.Equal(t, migration.Name, "add-uuid-extension")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

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
	migration := NewAddUuidExtensionMigration()
	assert.Equal(t, migration.Name, "add-uuid-extension")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WillDelayFor(time.Duration(50 * time.Millisecond)).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Millisecond))
	defer cancel()

	err = migration.Up(ctx, db, nil)
	assert.Nil(t, err)
}
