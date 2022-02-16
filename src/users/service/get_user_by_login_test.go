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
	"golang.org/x/crypto/bcrypt"
)

func Test_GetUserByLogin_FailIfInvalidInput(t *testing.T) {
	requests := []domain.GetUserRequest{
		{
			Username: "",
			Password: "",
		},
		{
			Username: "cenas",
			Password: "",
		},
		{
			Username: "",
			Password: "cenas",
		},
	}

	for _, req := range requests {
		service := newService(nil, nil)
		user, err := service.GetUserByLogin(context.TODO(), req)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid parameter")
	}
}

func Test_GetUserByLogin_FailIfGetByUsernameError(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUsername", mock.Anything, "username").Once().
		Return(nil, errors.New("boom"))

	service := newService(userRepo, nil)
	user, err := service.GetUserByLogin(context.TODO(), domain.GetUserRequest{Username: "username", Password: "casd"})
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
	userRepo.AssertExpectations(t)
}

func Test_GetUserByLogin_FailIfRoleFetchingError(t *testing.T) {
	roleUuid := uuid.New()

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUsername", mock.Anything, "username").Once().
		Return(&domain.User{
			ID:     uuid.New(),
			RoleId: roleUuid,
		}, nil)

	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetByUUID", mock.Anything, roleUuid).
		Once().Return(domain.Role{}, errors.New("boom"))

	service := newService(userRepo, roleRepo)
	user, err := service.GetUserByLogin(context.TODO(), domain.GetUserRequest{Username: "username", Password: "casd"})
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
	roleRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func Test_GetUserByLogin_FailIfPasswordsDontMatchError(t *testing.T) {
	roleUuid := uuid.New()

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUsername", mock.Anything, "username").Once().
		Return(&domain.User{
			ID:       uuid.New(),
			RoleId:   roleUuid,
			Password: "dont match password",
		}, nil)

	role := domain.Role{ID: roleUuid}
	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetByUUID", mock.Anything, roleUuid).
		Once().Return(role, nil)

	service := newService(userRepo, roleRepo)
	user, err := service.GetUserByLogin(context.TODO(), domain.GetUserRequest{Username: "username", Password: "casd"})
	assert.Nil(t, user)
	assert.Error(t, err)
	roleRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func Test_GetUserByLogin_Success(t *testing.T) {
	roleUuid := uuid.New()
	userUuid := uuid.New()

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte("password"), 14)
	assert.Nil(t, err)
	password := string(passwordBytes[:])

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUsername", mock.Anything, "username").Once().
		Return(&domain.User{
			ID:       userUuid,
			RoleId:   roleUuid,
			Password: password,
		}, nil)

	role := domain.Role{ID: roleUuid}
	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetByUUID", mock.Anything, roleUuid).
		Once().Return(role, nil)

	service := newService(userRepo, roleRepo)
	user, err := service.GetUserByLogin(context.TODO(), domain.GetUserRequest{Username: "username", Password: "password"})
	assert.Equal(t, user, &domain.User{
		ID:       userUuid,
		RoleId:   roleUuid,
		Password: password,
		Role:     &role,
	})
	assert.Nil(t, err)
	roleRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}
