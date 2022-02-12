package integration_tests

import (
	"context"
	"testing"

	users "github.com/plagioriginal/users-service-grpc/users"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Grpc_Add_User(t *testing.T) {
	tests := []struct {
		name             string
		req              *users.NewUserRequest
		loginReq         *users.LoginRequest
		wantedRes        *users.UserResponse
		wantedErrCode    codes.Code
		wantedErrMessage string
	}{
		{
			name:             "invalid request",
			loginReq:         nil,
			req:              &users.NewUserRequest{},
			wantedErrCode:    codes.Unauthenticated,
			wantedRes:        nil,
			wantedErrMessage: "invalid token",
		},
		{
			name: "invalid user role",
			loginReq: &users.LoginRequest{
				Username: databaseSettings.DefaultUserUsername,
				Password: databaseSettings.DefaultUserPassword,
			},
			req: &users.NewUserRequest{
				Username: "new-user",
				Password: "dummy-password",
				Role:     "invalid-user-role",
			},
			wantedRes:        nil,
			wantedErrCode:    codes.Internal,
			wantedErrMessage: "error storing user",
		},
		{
			name: "success",
			loginReq: &users.LoginRequest{
				Username: databaseSettings.DefaultUserUsername,
				Password: databaseSettings.DefaultUserPassword,
			},
			req: &users.NewUserRequest{
				Username: "new-user",
				Password: "dummy-password",
				Role:     "user",
			},
			wantedRes: &users.UserResponse{
				Username: "new-user",
				Role: &users.UserResponse_RoleResponse{
					RoleLabel: "User",
					RoleSlug:  "user",
				},
			},
			wantedErrCode:    0,
			wantedErrMessage: "",
		},
		{
			name: "user trying to add is not admin",
			loginReq: &users.LoginRequest{
				Username: "new-user",
				Password: "dummy-password",
			},
			req: &users.NewUserRequest{
				Username: "will not",
				Password: "be added",
				Role:     "user",
			},
			wantedRes:        nil,
			wantedErrCode:    codes.Unauthenticated,
			wantedErrMessage: "incorrect permissions",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			accessToken := ""

			if test.loginReq != nil {
				res, err := userClient.Login(context.Background(), test.loginReq)
				assert.Nil(t, err)
				assert.NotEmpty(t, res)
				accessToken = res.AccessToken
			}

			test.req.AccessToken = accessToken
			res, err := userClient.AddUser(context.TODO(), test.req)
			if test.wantedErrCode != 0 {
				assert.Equal(t, res, test.wantedRes)
				assert.Error(t, err)

				s, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, s.Code(), test.wantedErrCode)
				assert.Equal(t, s.Message(), test.wantedErrMessage)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, test.wantedRes.Username, res.Username)
			assert.Equal(t, test.wantedRes.Role.RoleLabel, res.Role.RoleLabel)
			assert.Equal(t, test.wantedRes.Role.RoleSlug, res.Role.RoleSlug)
		})
	}
}
