package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetByToken_ErrorPreparingContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	query := `
		SELECT id, token, valid_until
		FROM refresh_tokens
		WHERE token = $1
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).WillReturnError(errors.New("boom"))

	repo := PostgresRepository{db}
	res, err := repo.GetByToken(context.TODO(), uuid.New())
	assert.Equal(t, domain.RefreshToken{}, res)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func TestGetByToken_TimeoutReached(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	id := uuid.New()
	query := `
		SELECT id, token, valid_until
		FROM refresh_tokens
		WHERE token = $1
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(id).
		WillDelayFor(time.Duration(200 * time.Millisecond)).
		WillReturnError(errors.New("result doesnt matter because we are testing timeout"))

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(100*time.Millisecond))
	defer cancel()

	repo := PostgresRepository{db}
	res, err := repo.GetByToken(ctx, id)
	assert.Equal(t, domain.RefreshToken{}, res)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
}

func TestGetByToken_Success(t *testing.T) {
	id := uuid.New()
	token := uuid.New()
	createdAt := time.Now()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	expectedResult := sqlmock.NewRows(
		[]string{"id", "token", "valid_until"},
	).AddRow(id, token, createdAt)

	query := `
		SELECT id, token, valid_until
		FROM refresh_tokens
		WHERE token = $1
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(id).
		WillReturnRows(expectedResult)

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(100*time.Millisecond))
	defer cancel()

	repo := PostgresRepository{db}
	res, err := repo.GetByToken(ctx, id)
	assert.Equal(t, res.Id, id)
	assert.Equal(t, res.Token, token)
	assert.Equal(t, res.ValidUntil, createdAt)
	assert.Nil(t, err)
}
