package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/stretchr/testify/assert"
)

type anyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func Test_Store_ErrorPreparingContext(t *testing.T) {
	userId := uuid.New()
	roleId := uuid.New()
	user := domain.User{
		ID:        userId,
		FirstName: "cenas",
		Password:  "wrong password wtv",
		RoleId:    roleId,
	}

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	query := `INSERT INTO users(id, first_name, last_name, username, password, role_id, created_at, updated_at) 
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, first_name, last_name, username, password, role_id, refresh_token_id, created_at, updated_at`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).WillReturnError(errors.New("boom"))

	repo := PostgresRepository{db}
	createdUser, err := repo.Store(context.TODO(), user)
	assert.Nil(t, createdUser)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func Test_Store_TimeoutReached(t *testing.T) {
	userId := uuid.New()
	roleId := uuid.New()
	user := domain.User{
		ID:        userId,
		FirstName: "cenas",
		Password:  "wrong password wtv",
		RoleId:    roleId,
	}

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	query := `INSERT INTO users(id, first_name, last_name, username, password, role_id, created_at, updated_at) 
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, first_name, last_name, username, password, role_id, refresh_token_id, created_at, updated_at`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(
			userId,
			"cenas",
			"",
			"",
			"wrong password wtv",
			roleId,
			anyTime{},
			anyTime{},
		).WillDelayFor(time.Duration(6 * time.Second)).
		WillReturnError(errors.New("result doesnt matter because we are testing timeout"))

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(2*time.Second))
	defer cancel()

	repo := PostgresRepository{db}
	createdUser, err := repo.Store(ctx, user)
	assert.Nil(t, createdUser)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
}

func Test_Store_Success(t *testing.T) {
	userId := uuid.New()
	roleId := uuid.New()
	createdAt := time.Now()
	user := domain.User{
		ID:        userId,
		FirstName: "cenas",
		Password:  "wrong password wtv",
		RoleId:    roleId,
	}

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	expectedResult := sqlmock.NewRows(
		[]string{"id", "first_name", "last_name", "username", "password", "role_id", "refresh_token_id", "created_at", "updated_at"},
	).AddRow(userId, "cenas", "", "", "wrong password wtv", roleId, uuid.Nil, createdAt, createdAt)

	query := `INSERT INTO users(id, first_name, last_name, username, password, role_id, created_at, updated_at) 
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, first_name, last_name, username, password, role_id, refresh_token_id, created_at, updated_at`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(
			userId,
			"cenas",
			"",
			"",
			"wrong password wtv",
			roleId,
			anyTime{},
			anyTime{},
		).
		WillReturnRows(expectedResult)

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(2*time.Second))
	defer cancel()

	repo := PostgresRepository{db}
	createdUser, err := repo.Store(ctx, user)
	assert.Equal(t, createdUser.ID, userId)
	assert.Equal(t, createdUser.FirstName, "cenas")
	assert.Equal(t, createdUser.Password, "wrong password wtv")
	assert.Equal(t, createdUser.RoleId, roleId)
	assert.Equal(t, createdUser.CreatedAt, createdAt)
	assert.Equal(t, createdUser.UpdatedAt, createdAt)
	assert.Nil(t, err)
}
