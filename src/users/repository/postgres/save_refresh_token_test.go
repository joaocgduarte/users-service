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

func Test_SaveRefreshToken_ErrorPreparingContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	user := &domain.User{
		ID: uuid.New(),
	}
	token := domain.RefreshToken{
		Id: uuid.New(),
	}
	query := `
		UPDATE users
		SET refresh_token_id = $1
		WHERE id = $2
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).WillReturnError(errors.New("boom"))

	repo := PostgresRepository{db}
	err = repo.SaveRefreshToken(context.TODO(), user, token)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func Test_SaveRefreshToken_TimeoutReached(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	user := &domain.User{
		ID: uuid.New(),
	}
	token := domain.RefreshToken{
		Id: uuid.New(),
	}
	query := `
		UPDATE users
		SET refresh_token_id = $1
		WHERE id = $2
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectExec().
		WithArgs(token.Id, user.ID).
		WillDelayFor(time.Duration(200 * time.Millisecond)).
		WillReturnError(errors.New("result doesnt matter because we are testing timeout"))

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(100*time.Millisecond))
	defer cancel()

	repo := PostgresRepository{db}
	err = repo.SaveRefreshToken(ctx, user, token)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
}

func Test_SaveRefreshToken_ExecFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	user := &domain.User{
		ID: uuid.New(),
	}
	token := domain.RefreshToken{
		Id: uuid.New(),
	}
	query := `
		UPDATE users
		SET refresh_token_id = $1
		WHERE id = $2
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectExec().
		WithArgs(token.Id, user.ID).
		WillReturnError(errors.New("boom"))

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(100*time.Millisecond))
	defer cancel()

	repo := PostgresRepository{db}
	err = repo.SaveRefreshToken(ctx, user, token)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
}

func Test_SaveRefreshToken_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	user := &domain.User{
		ID: uuid.New(),
	}
	token := domain.RefreshToken{
		Id: uuid.New(),
	}
	query := `
		UPDATE users
		SET refresh_token_id = $1
		WHERE id = $2
	`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectExec().
		WithArgs(token.Id, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(100*time.Millisecond))
	defer cancel()

	repo := PostgresRepository{db}
	err = repo.SaveRefreshToken(ctx, user, token)
	assert.Nil(t, err)
}
