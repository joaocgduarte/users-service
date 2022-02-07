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

func TestCreateUsers_FailExec(t *testing.T) {
	migration := NewCreateUsersMigration()
	assert.Equal(t, migration.Name, "create-users-table")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

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

	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnError(errors.New("boom"))

	err = migration.Up(context.TODO(), db, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func TestCreateUsers_TimeoutReached(t *testing.T) {
	migration := NewCreateUsersMigration()
	assert.Equal(t, migration.Name, "create-users-table")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

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

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WillDelayFor(time.Duration(200 * time.Millisecond)).
		WillReturnError(errors.New("boom"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Millisecond))
	defer cancel()

	err = migration.Up(ctx, db, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
}

func TestCreateUsers_Success(t *testing.T) {
	migration := NewCreateUsersMigration()
	assert.Equal(t, migration.Name, "create-users-table")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

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

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WillDelayFor(time.Duration(50 * time.Millisecond)).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Millisecond))
	defer cancel()

	err = migration.Up(ctx, db, nil)
	assert.Nil(t, err)
}
