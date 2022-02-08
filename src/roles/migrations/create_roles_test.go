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

func TestCreateRoles_FailExec(t *testing.T) {
	migration := NewCreateRolesMigration()
	assert.Equal(t, migration.Name, "create-roles-table")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `CREATE TABLE IF NOT EXISTS roles(
		id uuid DEFAULT uuid_generate_v4(),
		role_slug varchar(255) NOT NULL UNIQUE,
		role_label varchar(255) NOT NULL UNIQUE,

		created_at timestamptz NOT NULL DEFAULT (now()),
		updated_at timestamptz NOT NULL DEFAULT (now()),
		deleted_at timestamptz DEFAULT NULL,
		PRIMARY KEY (id)
	);`

	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnError(errors.New("boom"))

	err = migration.Up(context.TODO(), db, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func TestCreateRoles_TimeoutReached(t *testing.T) {
	migration := NewCreateRolesMigration()
	assert.Equal(t, migration.Name, "create-roles-table")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `CREATE TABLE IF NOT EXISTS roles(
		id uuid DEFAULT uuid_generate_v4(),
		role_slug varchar(255) NOT NULL UNIQUE,
		role_label varchar(255) NOT NULL UNIQUE,

		created_at timestamptz NOT NULL DEFAULT (now()),
		updated_at timestamptz NOT NULL DEFAULT (now()),
		deleted_at timestamptz DEFAULT NULL,
		PRIMARY KEY (id)
	);`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WillDelayFor(time.Duration(200 * time.Millisecond)).
		WillReturnError(errors.New("boom"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Millisecond))
	defer cancel()

	err = migration.Up(ctx, db, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
}

func TestCreateRoles_Success(t *testing.T) {
	migration := NewCreateRolesMigration()
	assert.Equal(t, migration.Name, "create-roles-table")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `CREATE TABLE IF NOT EXISTS roles(
		id uuid DEFAULT uuid_generate_v4(),
		role_slug varchar(255) NOT NULL UNIQUE,
		role_label varchar(255) NOT NULL UNIQUE,

		created_at timestamptz NOT NULL DEFAULT (now()),
		updated_at timestamptz NOT NULL DEFAULT (now()),
		deleted_at timestamptz DEFAULT NULL,
		PRIMARY KEY (id)
	);`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WillDelayFor(time.Duration(50 * time.Millisecond)).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Millisecond))
	defer cancel()

	err = migration.Up(ctx, db, nil)
	assert.Nil(t, err)
}
