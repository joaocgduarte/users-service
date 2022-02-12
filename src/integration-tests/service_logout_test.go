package integration_tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	users "github.com/plagioriginal/users-service-grpc/users"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Grpc_Logout(t *testing.T) {
	tests := []struct {
		name             string
		loginReq         *users.LoginRequest
		wantedErrCode    codes.Code
		wantedErrMessage string
	}{
		{
			name:             "invalid request",
			loginReq:         nil,
			wantedErrCode:    codes.InvalidArgument,
			wantedErrMessage: "invalid input",
		},
		{
			name: "success",
			loginReq: &users.LoginRequest{
				Username: databaseSettings.DefaultUserUsername,
				Password: databaseSettings.DefaultUserPassword,
			},
			wantedErrCode:    0,
			wantedErrMessage: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := &users.RefreshRequest{}
			refreshToken := ""

			if test.loginReq != nil {
				res, err := userClient.Login(context.Background(), test.loginReq)
				assert.Nil(t, err)
				assert.NotEmpty(t, res)
				refreshToken = res.RefreshToken
			}

			req.RefreshToken = refreshToken
			res, err := userClient.Logout(context.Background(), req)
			if test.wantedErrCode != 0 {
				assert.Empty(t, res)
				assert.Error(t, err)

				s, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, s.Code(), test.wantedErrCode)
				assert.Equal(t, s.Message(), test.wantedErrMessage)
				return
			}

			// Result should always be empty access and refresh token
			assert.Nil(t, err)
			assert.Equal(t, res.AccessToken, "")
			assert.Equal(t, res.RefreshToken, "")

			// make sure old login token doesn't exist
			token, ok := uuid.Parse(refreshToken)
			assert.Nil(t, ok)
			assert.NotEqual(t, token, uuid.Nil)
			tokenRes, err := refreshTokenRepo.GetByToken(context.Background(), token)
			assert.Empty(t, tokenRes)
			assert.Error(t, err)
		})
	}
}
