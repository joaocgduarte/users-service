package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/domain/mocks"
	users "github.com/plagioriginal/users-service-grpc/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddUser_InvalidInput(t *testing.T) {
	service := newHandler(nil, nil)
	res, err := service.AddUser(context.TODO(), nil)
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.Unauthenticated, "invalid token"))

	res, err = service.AddUser(context.TODO(), &users.NewUserRequest{})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.Unauthenticated, "invalid token"))
}

func TestAddUser_ErrorParsingToken(t *testing.T) {
	accessTokenManager := new(mocks.AccessTokenHandler)
	service := newHandler(accessTokenManager, nil)

	accessTokenManager.On("ParseJWT", "cenas").Once().Return(nil, errors.New("boom"))
	res, err := service.AddUser(context.TODO(), &users.NewUserRequest{
		AccessToken: "cenas",
	})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.Unauthenticated, "invalid token"))
	accessTokenManager.AssertExpectations(t)
}

func TestAddUser_ErrorIfTokenIsInvalid(t *testing.T) {
	accessTokenManager := new(mocks.AccessTokenHandler)
	service := newHandler(accessTokenManager, nil)

	mockToken := &jwt.Token{Raw: "mock token"}
	accessTokenManager.On("ParseJWT", "cenas").Once().Return(mockToken, nil)
	accessTokenManager.On("IsJWTokenValid", mockToken).Once().Return(false)
	res, err := service.AddUser(context.TODO(), &users.NewUserRequest{
		AccessToken: "cenas",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.Unauthenticated, "invalid token"))
	accessTokenManager.AssertExpectations(t)
}

func TestAddUser_ErrorGettingUserFromToken(t *testing.T) {
	accessTokenManager := new(mocks.AccessTokenHandler)
	service := newHandler(accessTokenManager, nil)

	mockToken := &jwt.Token{Raw: "mock token"}
	accessTokenManager.On("ParseJWT", "cenas").Once().Return(mockToken, nil)
	accessTokenManager.On("IsJWTokenValid", mockToken).Once().Return(true)
	accessTokenManager.On("GetUserRoleFromToken", mockToken).Once().Return("", errors.New("boom"))
	res, err := service.AddUser(context.TODO(), &users.NewUserRequest{
		AccessToken: "cenas",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.InvalidArgument, "invalid token"))
	accessTokenManager.AssertExpectations(t)
}

func TestAddUser_UserDoesntHaveProperRole(t *testing.T) {
	accessTokenManager := new(mocks.AccessTokenHandler)
	service := newHandler(accessTokenManager, nil)

	mockToken := &jwt.Token{Raw: "mock token"}
	accessTokenManager.On("ParseJWT", "cenas").Once().Return(mockToken, nil)
	accessTokenManager.On("IsJWTokenValid", mockToken).Once().Return(true)
	accessTokenManager.On("GetUserRoleFromToken", mockToken).Once().Return("user", nil)
	res, err := service.AddUser(context.TODO(), &users.NewUserRequest{
		AccessToken: "cenas",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.Unauthenticated, "incorrect permissions"))
	accessTokenManager.AssertExpectations(t)
}

func TestAddUser_InvalidRequestInParameters(t *testing.T) {
	accessTokenManager := new(mocks.AccessTokenHandler)
	service := newHandler(accessTokenManager, nil)

	mockToken := &jwt.Token{Raw: "mock token"}
	accessTokenManager.On("ParseJWT", "cenas").Once().Return(mockToken, nil)
	accessTokenManager.On("IsJWTokenValid", mockToken).Once().Return(true)
	accessTokenManager.On("GetUserRoleFromToken", mockToken).Once().Return("admin", nil)
	res, err := service.AddUser(context.TODO(), &users.NewUserRequest{
		AccessToken: "cenas",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	accessTokenManager.AssertExpectations(t)
}

func TestAddUser_ErrorUserStoring(t *testing.T) {
	accessTokenManager := new(mocks.AccessTokenHandler)
	userService := new(mocks.UserService)
	service := newHandler(accessTokenManager, userService)

	mockToken := &jwt.Token{Raw: "mock token"}
	accessTokenManager.On("ParseJWT", "cenas").Once().Return(mockToken, nil)
	accessTokenManager.On("IsJWTokenValid", mockToken).Once().Return(true)
	accessTokenManager.On("GetUserRoleFromToken", mockToken).Once().Return("admin", nil)

	userService.On("Store", mock.Anything, domain.StoreUserRequest{
		Username: "username",
		Password: "password",
		RoleSlug: "some role",
	}).Once().Return(nil, errors.New("boom"))

	res, err := service.AddUser(context.TODO(), &users.NewUserRequest{
		AccessToken: "cenas",
		Username:    "username",
		Password:    "password",
		Role:        "some role",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.Unknown, "error storing user"))
	accessTokenManager.AssertExpectations(t)
	userService.AssertExpectations(t)
}

func TestAddUser_Success(t *testing.T) {
	accessTokenManager := new(mocks.AccessTokenHandler)
	userService := new(mocks.UserService)
	service := newHandler(accessTokenManager, userService)

	mockToken := &jwt.Token{Raw: "mock token"}
	accessTokenManager.On("ParseJWT", "cenas").Once().Return(mockToken, nil)
	accessTokenManager.On("IsJWTokenValid", mockToken).Once().Return(true)
	accessTokenManager.On("GetUserRoleFromToken", mockToken).Once().Return("admin", nil)

	userId := uuid.New()
	roleId := uuid.New()

	userService.On("Store", mock.Anything, domain.StoreUserRequest{
		Username: "username",
		Password: "password",
		RoleSlug: "some role",
	}).Once().Return(&domain.User{
		ID:       userId,
		Username: "username",
		Password: "encrypted password",
		RoleId:   roleId,
		Role: &domain.Role{
			ID:        roleId,
			RoleSlug:  "admin",
			RoleLabel: "Administrator",
		},
	}, nil)

	res, err := service.AddUser(context.TODO(), &users.NewUserRequest{
		AccessToken: "cenas",
		Username:    "username",
		Password:    "password",
		Role:        "some role",
	})

	assert.Equal(t, res, &users.UserResponse{
		Id:        userId.String(),
		Username:  "username",
		FirstName: "",
		LastName:  "",
		Role: &users.UserResponse_RoleResponse{
			Id:        roleId.String(),
			RoleLabel: "Administrator",
			RoleSlug:  "admin",
		},
	})
	assert.Nil(t, err)
	accessTokenManager.AssertExpectations(t)
	userService.AssertExpectations(t)
}
