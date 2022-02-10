package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTokenFromRepo_ErrorIfRepoFails(t *testing.T) {
	token := uuid.New()

	tokenRepo := new(mocks.RefreshTokenRepository)
	tokenRepo.On("GetByToken", mock.Anything, token).
		Once().Return(domain.RefreshToken{}, errors.New("boom"))

	res, err := newService(tokenRepo, nil).
		GetTokenFromRepo(context.TODO(), token)

	assert.Empty(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
	tokenRepo.AssertExpectations(t)
}

func TestGetTokenFromRepo_Success(t *testing.T) {
	token := uuid.New()

	tokenRepo := new(mocks.RefreshTokenRepository)
	tokenRepo.On("GetByToken", mock.Anything, token).
		Once().Return(domain.RefreshToken{
		Token: token,
	}, nil)

	res, err := newService(tokenRepo, nil).
		GetTokenFromRepo(context.TODO(), token)

	assert.Nil(t, err)
	assert.Equal(t, res, domain.RefreshToken{
		Token: token,
	})
	tokenRepo.AssertExpectations(t)
}
