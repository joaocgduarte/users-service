package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/domain/mocks"
	users "github.com/plagioriginal/users-service-grpc/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogin_InvalidInput(t *testing.T) {
	service := newHandler(nil, nil)
	res, err := service.Login(context.TODO(), nil)
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.NotFound, "resource found"))

	res, err = service.Login(context.TODO(), &users.LoginRequest{})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.NotFound, "resource found"))
}

func TestLogin_ErrorGettingTheUserByLogin(t *testing.T) {
	userService := new(mocks.UserService)
	service := newHandler(nil, userService)

	userService.On("GetUserByLogin", mock.Anything, domain.GetUserRequest{
		Username: "username",
		Password: "password",
	}).Once().Return(nil, errors.New("boom"))

	res, err := service.Login(context.TODO(), &users.LoginRequest{
		Username: "username",
		Password: "password",
	})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.NotFound, "resource found"))
	userService.AssertExpectations(t)
}

func TestLogin_ErrorGeneratingToken(t *testing.T) {
	userService := new(mocks.UserService)
	tokenHandler := new(mocks.AccessTokenHandler)
	service := newHandler(tokenHandler, userService)

	userService.On("GetUserByLogin", mock.Anything, domain.GetUserRequest{
		Username: "username",
		Password: "password",
	}).Once().Return(&domain.User{}, nil)

	tokenHandler.On("GenerateTokens", mock.Anything, &domain.User{}).Once().
		Return(domain.TokenResponse{}, errors.New("boom"))

	res, err := service.Login(context.TODO(), &users.LoginRequest{
		Username: "username",
		Password: "password",
	})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.Unknown, "error generating tokens"))
	userService.AssertExpectations(t)
	tokenHandler.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	userService := new(mocks.UserService)
	tokenHandler := new(mocks.AccessTokenHandler)
	service := newHandler(tokenHandler, userService)

	userId := uuid.New()
	roleId := uuid.New()

	user := &domain.User{
		ID:       userId,
		Username: "username",
		Password: "encrypted password",
		RoleId:   roleId,
		Role: &domain.Role{
			ID:        roleId,
			RoleSlug:  "admin",
			RoleLabel: "Administrator",
		},
	}
	userService.On("GetUserByLogin", mock.Anything, domain.GetUserRequest{
		Username: "username",
		Password: "password",
	}).Once().Return(user, nil)

	tokenHandler.On("GenerateTokens", mock.Anything, user).Once().
		Return(domain.TokenResponse{AccessToken: "access-token", RefreshToken: "refresh-token"}, nil)

	res, err := service.Login(context.TODO(), &users.LoginRequest{
		Username: "username",
		Password: "password",
	})
	assert.Equal(t, res, &users.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		User: &users.UserResponse{
			Id:        userId.String(),
			Username:  "username",
			FirstName: "",
			LastName:  "",
			Role: &users.UserResponse_RoleResponse{
				Id:        roleId.String(),
				RoleLabel: "Administrator",
				RoleSlug:  "admin",
			},
		},
	})
	assert.Nil(t, err)
	userService.AssertExpectations(t)
	tokenHandler.AssertExpectations(t)
}
