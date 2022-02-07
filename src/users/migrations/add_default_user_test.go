package migrations

import (
	"context"
	"io/ioutil"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Rest of the cases are already tested in service store test

func TestAddDefaultUserMigration_FailtIfNoDefaultUserInContext(t *testing.T) {
	migration := NewAddDefaultUserMigration()
	assert.Equal(t, migration.Name, "add-default-user")

	err := migration.Up(context.TODO(), nil, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "cant add default users without credentials")
}

func TestAddDefaultUserMigration_QueryError(t *testing.T) {
	migration := NewAddDefaultUserMigration()
	assert.Equal(t, migration.Name, "add-default-user")

	db, _, err := sqlmock.New()
	assert.Nil(t, err)

	ctx := context.WithValue(context.TODO(), "defaultUserUsername", "admin")
	ctx = context.WithValue(ctx, "defaultUserPassword", "admin")

	err = migration.Up(ctx, db, log.New(ioutil.Discard, "tests: ", log.Flags()))
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error fetching role")
}
