package integration_tests

import (
	"context"
	"testing"

	users "github.com/plagioriginal/users-service-grpc/users"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Grpc_Login(t *testing.T) {
	tests := []struct {
		name             string
		req              *users.LoginRequest
		wantedRes        *users.TokenResponse
		wantedErrCode    codes.Code
		wantedErrMessage string
	}{
		{
			name:             "nil login request",
			req:              &users.LoginRequest{},
			wantedRes:        nil,
			wantedErrCode:    codes.NotFound,
			wantedErrMessage: "resource found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := userClient.Login(context.Background(), test.req)

			if test.wantedErrCode != 0 {
				s, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, s.Message(), test.wantedErrMessage)
				assert.Equal(t, test.wantedRes, res)
				assert.Equal(t, s.Code(), test.wantedErrCode)
				return
			}

			assert.Nil(t, err)
			assert.EqualValues(t, test.wantedRes, res)
		})
	}
	userClient.Login(context.Background(), nil)
}
