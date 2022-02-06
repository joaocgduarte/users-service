package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetByUsername_ErrorPreparingContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	query := `
		SELECT id, first_name, last_name, username, password, role_id, refresh_token_id, created_at, updated_at
		FROM users 
		WHERE username = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).WillReturnError(errors.New("boom"))

	repo := PostgresRepository{db}
	user, err := repo.GetByUsername(context.TODO(), "cenas")
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func Test_GetByUsername_TimeoutReached(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	query := `
		SELECT id, first_name, last_name, username, password, role_id, refresh_token_id, created_at, updated_at
		FROM users 
		WHERE username = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs("admin").
		WillDelayFor(time.Duration(3 * time.Second)).
		WillReturnError(errors.New("result doesnt matter because we are testing timeout"))

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(2*time.Second))
	defer cancel()

	repo := PostgresRepository{db}
	user, err := repo.GetByUsername(ctx, "admin")
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
}

func Test_GetByUsername_Success(t *testing.T) {
	userId := uuid.New()
	roleId := uuid.New()
	createdAt := time.Now()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	expectedResult := sqlmock.NewRows(
		[]string{"id", "first_name", "last_name", "username", "password", "role_id", "refresh_token_id", "created_at", "updated_at"},
	).AddRow(userId, "cenas", "", "admin", "wrong password wtv", roleId, uuid.Nil, createdAt, createdAt)

	query := `
		SELECT id, first_name, last_name, username, password, role_id, refresh_token_id, created_at, updated_at
		FROM users 
		WHERE username = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs("admin").
		WillReturnRows(expectedResult)

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(2*time.Second))
	defer cancel()

	repo := PostgresRepository{db}
	createdUser, err := repo.GetByUsername(ctx, "admin")
	assert.Equal(t, createdUser.ID, userId)
	assert.Equal(t, createdUser.FirstName, "cenas")
	assert.Equal(t, createdUser.Username, "admin")
	assert.Equal(t, createdUser.Password, "wrong password wtv")
	assert.Equal(t, createdUser.RoleId, roleId)
	assert.Equal(t, createdUser.CreatedAt, createdAt)
	assert.Equal(t, createdUser.UpdatedAt, createdAt)
	assert.Nil(t, err)
}
