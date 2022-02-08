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

func TestStore_FailPrepareQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	id := uuid.New()
	token := uuid.New()
	validUntil := time.Now()

	query := `
		INSERT INTO refresh_tokens(id, token, valid_until)
		VALUES ($1, $2, $3)
		RETURNING id, token, valid_until
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(id, token, validUntil).
		WillReturnError(errors.New("boom"))

	res, err := New(db).Store(context.TODO(), domain.RefreshToken{
		Id:         id,
		Token:      token,
		ValidUntil: validUntil,
	})

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
	assert.Empty(t, res)
}

func TestStore_TimeoutReached(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	id := uuid.New()
	token := uuid.New()
	validUntil := time.Now()

	query := `
		INSERT INTO refresh_tokens(id, token, valid_until)
		VALUES ($1, $2, $3)
		RETURNING id, token, valid_until
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(id, token, validUntil).
		WillDelayFor(time.Duration(time.Millisecond * 150)).
		WillReturnError(errors.New("doessn't matter"))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Millisecond*100))
	defer cancel()

	res, err := New(db).Store(ctx, domain.RefreshToken{
		Id:         id,
		Token:      token,
		ValidUntil: validUntil,
	})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
	assert.Empty(t, res)
}

func TestStore_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	id := uuid.New()
	token := uuid.New()
	validUntil := time.Now()

	query := `
		INSERT INTO refresh_tokens(id, token, valid_until)
		VALUES ($1, $2, $3)
		RETURNING id, token, valid_until
	`

	expectedResult := sqlmock.NewRows([]string{"id", "token", "valid_until"})
	expectedResult.AddRow(id, token, validUntil)

	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(id, token, validUntil).
		WillReturnError(nil).
		WillReturnRows(expectedResult)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Millisecond*150))
	defer cancel()

	res, err := New(db).Store(ctx, domain.RefreshToken{
		Id:         id,
		Token:      token,
		ValidUntil: validUntil,
	})
	assert.Nil(t, err)
	assert.Equal(t, res, domain.RefreshToken{
		Id:         id,
		Token:      token,
		ValidUntil: validUntil,
	})
}
