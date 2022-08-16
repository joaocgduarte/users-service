package integration_tests

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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
		{
			name:             "nil login request",
			req:              &users.LoginRequest{},
			wantedRes:        nil,
			wantedErrCode:    codes.NotFound,
			wantedErrMessage: "resource found",
		},
		{
			name: "login with invalid username",
			req: &users.LoginRequest{
				Username: "invalid username",
				Password: "doesn't matter",
			},
			wantedRes:        nil,
			wantedErrCode:    codes.NotFound,
			wantedErrMessage: "resource found",
		},
		{
			name: "login with invalid password",
			req: &users.LoginRequest{
				Username: "default-user",
				Password: "doesn't matter",
			},
			wantedRes:        nil,
			wantedErrCode:    codes.NotFound,
			wantedErrMessage: "resource found",
		},
		{
			name: "success",
			req: &users.LoginRequest{
				Username: databaseSettings.DefaultUserUsername,
				Password: databaseSettings.DefaultUserPassword,
			},
			wantedRes: &users.TokenResponse{
				User: &users.UserResponse{
					Username: databaseSettings.DefaultUserUsername,
					Role: &users.UserResponse_RoleResponse{
						RoleLabel: "Administrator",
						RoleSlug:  "admin",
					},
				},
			},
			wantedErrCode:    0,
			wantedErrMessage: "",
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
			assert.Equal(t, test.wantedRes.User.Username, res.User.Username)
			assert.Equal(t, test.wantedRes.User.Role.RoleLabel, res.User.Role.RoleLabel)
			assert.Equal(t, test.wantedRes.User.Role.RoleSlug, res.User.Role.RoleSlug)

			accessToken := res.AccessToken
			refreshToken := res.RefreshToken

			jwt, err := jwt.Parse(accessToken, func(*jwt.Token) (interface{}, error) {
				return []byte(databaseSettings.JwtSecret), nil
			})
			assert.Nil(t, err)
			assert.True(t, jwt.Valid)

			refreshTokenUuid, err := uuid.Parse(refreshToken)
			assert.Nil(t, err)
			assert.NotEqual(t, refreshTokenUuid, uuid.Nil)

		})
	}
	userClient.Login(context.Background(), nil)
}
