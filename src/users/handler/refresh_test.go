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

func TestRefresh_InvalidInput(t *testing.T) {
	service := newHandler(nil, nil)
	res, err := service.Refresh(context.TODO(), nil)
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.NotFound, "resource found"))

	res, err = service.Refresh(context.TODO(), &users.RefreshRequest{})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.NotFound, "resource found"))
}

func TestRefresh_ErrorParsingOldToken(t *testing.T) {
	service := newHandler(nil, nil)
	res, err := service.Refresh(context.TODO(), &users.RefreshRequest{
		RefreshToken: "cenas",
	})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.InvalidArgument, "invalid token"))
}

func TestRefresh_ErrorRefreshingTokens(t *testing.T) {
	tokenHandler := new(mocks.AccessTokenHandler)
	service := newHandler(tokenHandler, nil)

	oldRefreshToken, err := uuid.Parse("a20b5aec-7000-4828-ad56-9d30675a49f2")
	assert.Nil(t, err)

	tokenHandler.On("RefreshAllTokens", mock.Anything, oldRefreshToken).
		Once().
		Return(domain.TokenResponse{}, errors.New("boom"))

	res, err := service.Refresh(context.TODO(), &users.RefreshRequest{
		RefreshToken: "a20b5aec-7000-4828-ad56-9d30675a49f2",
	})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.Internal, "error generating tokens"))
	tokenHandler.AssertExpectations(t)
}

func TestRefresh_ErrorParsingAccessToken(t *testing.T) {
	tokenHandler := new(mocks.AccessTokenHandler)
	service := newHandler(tokenHandler, nil)

	oldRefreshToken, err := uuid.Parse("a20b5aec-7000-4828-ad56-9d30675a49f2")
	assert.Nil(t, err)

	tokenHandler.On("RefreshAllTokens", mock.Anything, oldRefreshToken).
		Once().
		Return(domain.TokenResponse{
			AccessToken:  "cenas",
			RefreshToken: "refresh-cenas",
		}, nil)
	tokenHandler.On("ParseJWT", "cenas").
		Once().
		Return(nil, errors.New("boom"))

	res, err := service.Refresh(context.TODO(), &users.RefreshRequest{
		RefreshToken: "a20b5aec-7000-4828-ad56-9d30675a49f2",
	})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.Internal, "error generating tokens"))
	tokenHandler.AssertExpectations(t)
}

func TestRefresh_ErrorGettingUserIdFromToken(t *testing.T) {
	tokenHandler := new(mocks.AccessTokenHandler)
	service := newHandler(tokenHandler, nil)

	oldRefreshToken, err := uuid.Parse("a20b5aec-7000-4828-ad56-9d30675a49f2")
	assert.Nil(t, err)

	tokenHandler.On("RefreshAllTokens", mock.Anything, oldRefreshToken).
		Once().
		Return(domain.TokenResponse{
			AccessToken:  "cenas",
			RefreshToken: "refresh-cenas",
		}, nil)
	tokenHandler.On("ParseJWT", "cenas").
		Once().
		Return(&jwt.Token{Raw: "cenas"}, nil)
	tokenHandler.On("GetUserIDFromToken", &jwt.Token{Raw: "cenas"}).
		Once().Return(uuid.Nil, errors.New("boom"))

	res, err := service.Refresh(context.TODO(), &users.RefreshRequest{
		RefreshToken: "a20b5aec-7000-4828-ad56-9d30675a49f2",
	})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.Internal, "error generating tokens"))
	tokenHandler.AssertExpectations(t)
}

func TestRefresh_ErrorGettingUserByUuid(t *testing.T) {
	tokenHandler := new(mocks.AccessTokenHandler)
	userService := new(mocks.UserService)
	service := newHandler(tokenHandler, userService)

	oldRefreshToken, err := uuid.Parse("a20b5aec-7000-4828-ad56-9d30675a49f2")
	assert.Nil(t, err)

	userId := uuid.New()
	tokenHandler.On("RefreshAllTokens", mock.Anything, oldRefreshToken).
		Once().
		Return(domain.TokenResponse{
			AccessToken:  "cenas",
			RefreshToken: "refresh-cenas",
		}, nil)
	tokenHandler.On("ParseJWT", "cenas").
		Once().
		Return(&jwt.Token{Raw: "cenas"}, nil)
	tokenHandler.On("GetUserIDFromToken", &jwt.Token{Raw: "cenas"}).
		Once().Return(userId, nil)

	userService.On("GetUserByUUID", mock.Anything, userId).Once().Return(nil, errors.New("boom"))

	res, err := service.Refresh(context.TODO(), &users.RefreshRequest{
		RefreshToken: "a20b5aec-7000-4828-ad56-9d30675a49f2",
	})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.Internal, "error generating tokens"))
	tokenHandler.AssertExpectations(t)
	userService.AssertExpectations(t)
}

func TestRefresh_Success(t *testing.T) {
	tokenHandler := new(mocks.AccessTokenHandler)
	userService := new(mocks.UserService)
	service := newHandler(tokenHandler, userService)

	oldRefreshToken, err := uuid.Parse("a20b5aec-7000-4828-ad56-9d30675a49f2")
	assert.Nil(t, err)

	userId := uuid.New()
	roleId := uuid.New()

	tokenHandler.On("RefreshAllTokens", mock.Anything, oldRefreshToken).
		Once().
		Return(domain.TokenResponse{
			AccessToken:  "cenas",
			RefreshToken: "refresh-cenas",
		}, nil)
	tokenHandler.On("ParseJWT", "cenas").
		Once().
		Return(&jwt.Token{Raw: "cenas"}, nil)
	tokenHandler.On("GetUserIDFromToken", &jwt.Token{Raw: "cenas"}).
		Once().Return(userId, nil)

	userService.On("GetUserByUUID", mock.Anything, userId).Once().
		Return(&domain.User{
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

	res, err := service.Refresh(context.TODO(), &users.RefreshRequest{
		RefreshToken: "a20b5aec-7000-4828-ad56-9d30675a49f2",
	})
	assert.Equal(t, res, &users.TokenResponse{
		AccessToken:  "cenas",
		RefreshToken: "refresh-cenas",
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
	tokenHandler.AssertExpectations(t)
	userService.AssertExpectations(t)
}
