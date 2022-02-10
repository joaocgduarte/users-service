package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserByToken_ErrorIfRepoFails(t *testing.T) {
	token := domain.RefreshToken{
		Id:         uuid.New(),
		Token:      uuid.New(),
		ValidUntil: time.Now().Add(time.Duration(5) * time.Hour),
	}

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByRefreshToken", mock.Anything, token.Id).
		Once().Return(nil, errors.New("boom"))

	res, err := newService(nil, userRepo).
		GetUserByToken(context.TODO(), token)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
	userRepo.AssertExpectations(t)
}

func TestGetUserByToken_Success(t *testing.T) {
	token := domain.RefreshToken{
		Id:         uuid.New(),
		Token:      uuid.New(),
		ValidUntil: time.Now().Add(time.Duration(5) * time.Hour),
	}

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByRefreshToken", mock.Anything, token.Id).
		Once().Return(&domain.User{
		Username: "cenas",
	}, nil)

	res, err := newService(nil, userRepo).
		GetUserByToken(context.TODO(), token)

	assert.Nil(t, err)
	assert.Equal(t, res, &domain.User{
		Username: "cenas",
	})
	userRepo.AssertExpectations(t)
}
