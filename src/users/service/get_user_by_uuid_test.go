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

func Test_GetUserUUID_FailIfGetByUsernameError(t *testing.T) {
	uuid := uuid.New()
	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUUID", mock.Anything, uuid).Once().
		Return(nil, errors.New("boom"))

	service := newService(userRepo, nil)
	user, err := service.GetUserByUUID(context.TODO(), uuid)
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")
	userRepo.AssertExpectations(t)
}

func Test_GetUserUUID_FailIfRoleFetchingError(t *testing.T) {
	userUuid := uuid.New()
	roleUuid := uuid.New()

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUUID", mock.Anything, userUuid).Once().
		Return(&domain.User{
			ID:     userUuid,
			RoleId: roleUuid,
		}, nil)

	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetByUUID", mock.Anything, roleUuid).
		Once().Return(domain.Role{}, errors.New("boom"))

	service := newService(userRepo, roleRepo)
	user, err := service.GetUserByUUID(context.TODO(), userUuid)
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")
	roleRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func Test_GetUserUUID_Success(t *testing.T) {
	userUuid := uuid.New()
	roleUuid := uuid.New()

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUUID", mock.Anything, userUuid).Once().
		Return(&domain.User{
			ID:     userUuid,
			RoleId: roleUuid,
		}, nil)

	role := domain.Role{ID: roleUuid}
	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetByUUID", mock.Anything, roleUuid).
		Once().Return(role, nil)

	service := newService(userRepo, roleRepo)
	user, err := service.GetUserByUUID(context.TODO(), userUuid)
	assert.Equal(t, user, &domain.User{
		ID:     userUuid,
		RoleId: roleUuid,
		Role:   &role,
	})
	assert.Nil(t, err)
	roleRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}
