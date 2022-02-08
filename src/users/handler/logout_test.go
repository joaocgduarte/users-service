package handler

import (
	"context"
	"testing"

	"github.com/plagioriginal/user-microservice/domain/mocks"
	users "github.com/plagioriginal/users-service-grpc/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogout_InvalidInput(t *testing.T) {
	service := newHandler(nil, nil)
	res, err := service.Logout(context.TODO(), nil)
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.InvalidArgument, "invalid input"))

	res, err = service.Logout(context.TODO(), &users.RefreshRequest{})
	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, status.Error(codes.InvalidArgument, "invalid input"))
}

func TestLogout_Success(t *testing.T) {
	tokenManager := new(mocks.AccessTokenHandler)
	service := newHandler(tokenManager, nil)

	tokenManager.On("DeleteRefreshToken", mock.Anything, "cenas").Once().Return(true)
	res, err := service.Logout(context.TODO(), &users.RefreshRequest{RefreshToken: "cenas"})
	assert.Nil(t, err)
	assert.Equal(t, res, &users.TokenResponse{
		AccessToken:  "",
		RefreshToken: "",
	})
}
