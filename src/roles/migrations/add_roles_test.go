package migrations

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddDefaultRoles_FailExec(t *testing.T) {
	migration := NewAddRolesMigration()
	assert.Equal(t, migration.Name, "add-default-roles")

	db, _, err := sqlmock.New()
	assert.Nil(t, err)

	err = migration.Up(context.TODO(), db, log.New(ioutil.Discard, "tests: ", log.Flags()))
	assert.Error(t, err)
}

func TestAddDefaultRoles_SkipIfRoleExists(t *testing.T) {
	migration := NewAddRolesMigration()
	assert.Equal(t, migration.Name, "add-default-roles")

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	createdAt := time.Now()
	expectedResult := sqlmock.NewRows([]string{"id", "role_slug", "role_label", "created_at", "updated_at"})
	expectedResult.AddRow(uuid.New(), "admin", "Administrator", createdAt, createdAt)
	expectedResult2 := sqlmock.NewRows([]string{"id", "role_slug", "role_label", "created_at", "updated_at"})

	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL and role_slug=$1;`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs("admin").
		WillReturnRows(expectedResult).
		WillReturnError(nil)

	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs("user").
		WillReturnRows(expectedResult2).
		WillReturnError(errors.New("not existant"))

	insertedResult := sqlmock.NewRows([]string{"id", "role_slug", "role_label", "created_at", "updated_at"})
	insertedResult.AddRow(uuid.New(), "user", "User", createdAt, createdAt)

	query = `
		INSERT INTO roles (id, role_slug, role_label, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, role_slug, role_label, created_at, updated_at
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(sqlmock.AnyArg(), "user", "User", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(insertedResult).
		WillReturnError(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Millisecond))
	defer cancel()

	err = migration.Up(ctx, db, log.New(ioutil.Discard, "tests: ", log.Flags()))
	assert.Nil(t, err)
}
