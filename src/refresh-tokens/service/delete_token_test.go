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

func TestDeleteToken_ErrorIfRepoFails(t *testing.T) {
	id := uuid.New()

	tokenRepo := new(mocks.RefreshTokenRepository)
	tokenRepo.On("Delete", mock.Anything, id).
		Once().Return(errors.New("boom"))

	err := newService(tokenRepo, nil).
		DeleteToken(context.TODO(), domain.RefreshToken{
			Id: id,
		})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
	tokenRepo.AssertExpectations(t)
}

func TestDeleteToken_Success(t *testing.T) {
	id := uuid.New()

	tokenRepo := new(mocks.RefreshTokenRepository)
	tokenRepo.On("Delete", mock.Anything, id).
		Once().Return(nil)

	err := newService(tokenRepo, nil).
		DeleteToken(context.TODO(), domain.RefreshToken{
			Id: id,
		})

	assert.Nil(t, err)
	tokenRepo.AssertExpectations(t)
}
