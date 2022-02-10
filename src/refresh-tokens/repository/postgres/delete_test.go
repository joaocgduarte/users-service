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

func TestDelete_ErrorPreparingContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	query := `
		DELETE FROM refresh_tokens
		WHERE id = $1
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).WillReturnError(errors.New("boom"))

	repo := PostgresRepository{db}
	err = repo.Delete(context.TODO(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func TestDelete_TimeoutReached(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	id := uuid.New()
	query := `
		DELETE FROM refresh_tokens
		WHERE id = $1
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectExec().
		WithArgs(id).
		WillDelayFor(time.Duration(200 * time.Millisecond)).
		WillReturnError(errors.New("result doesnt matter because we are testing timeout"))

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(100*time.Millisecond))
	defer cancel()

	repo := PostgresRepository{db}
	err = repo.Delete(ctx, id)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
}

func TestDelete_ExecFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	id := uuid.New()
	query := `
		DELETE FROM refresh_tokens
		WHERE id = $1
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectExec().
		WithArgs(id).
		WillReturnError(errors.New("boom"))

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(100*time.Millisecond))
	defer cancel()

	repo := PostgresRepository{db}
	err = repo.Delete(ctx, id)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	id := uuid.New()
	query := `
		DELETE FROM refresh_tokens
		WHERE id = $1
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectExec().
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(100*time.Millisecond))
	defer cancel()

	repo := PostgresRepository{db}
	err = repo.Delete(ctx, id)
	assert.Nil(t, err)
}
