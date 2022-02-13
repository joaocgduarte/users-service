package integration_tests

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	users "github.com/plagioriginal/users-service-grpc/users"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Grpc_Refresh(t *testing.T) {
	tests := []struct {
		name                  string
		loginReq              *users.LoginRequest
		refreshTokenOverwrite string
		wantedRes             *users.TokenResponse
		wantedErrCode         codes.Code
		wantedErrMessage      string
	}{
		{
			name:                  "invalid request",
			loginReq:              nil,
			refreshTokenOverwrite: "",
			wantedRes:             nil,
			wantedErrCode:         codes.NotFound,
			wantedErrMessage:      "resource found",
		},
		{
			name:                  "invalid refresh-token",
			loginReq:              nil,
			refreshTokenOverwrite: "bla-blablab-abla-ali-g",
			wantedRes:             nil,
			wantedErrCode:         codes.InvalidArgument,
			wantedErrMessage:      "invalid token",
		},
		{
			name:                  "not existant refresh-token",
			loginReq:              nil,
			refreshTokenOverwrite: uuid.New().String(),
			wantedRes:             nil,
			wantedErrCode:         codes.Internal,
			wantedErrMessage:      "error generating tokens",
		},
		{
			name: "success",
			loginReq: &users.LoginRequest{
				Username: databaseSettings.DefaultUserUsername,
				Password: databaseSettings.DefaultUserPassword,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := &users.RefreshRequest{}
			refreshToken := ""
			var oldRefreshToken *domain.RefreshToken

			if test.loginReq != nil {
				res, err := userClient.Login(context.Background(), test.loginReq)
				assert.Nil(t, err)
				assert.NotEmpty(t, res)
				refreshToken = res.RefreshToken

				// make sure the login token exists
				token, ok := uuid.Parse(refreshToken)
				assert.Nil(t, ok)
				assert.NotEqual(t, token, uuid.Nil)
				tokenRes, err := refreshTokenRepo.GetByToken(context.Background(), token)
				assert.Equal(t, tokenRes.Token.String(), refreshToken)
				assert.Nil(t, err)
				oldRefreshToken = &tokenRes
			}

			req.RefreshToken = refreshToken
			if test.refreshTokenOverwrite != "" {
				req.RefreshToken = test.refreshTokenOverwrite
			}

			res, err := userClient.Refresh(context.Background(), req)
			if test.wantedErrCode != 0 {
				assert.Error(t, err)
				assert.Nil(t, res)

				s, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, s.Code(), test.wantedErrCode)
				assert.Equal(t, s.Message(), test.wantedErrMessage)
				return
			}

			accessToken := res.AccessToken
			newRefreshToken := res.RefreshToken

			jwt, err := jwt.Parse(accessToken, func(*jwt.Token) (interface{}, error) {
				return []byte(databaseSettings.JwtSecret), nil
			})
			assert.Nil(t, err)
			assert.True(t, jwt.Valid)

			refreshTokenUuid, err := uuid.Parse(newRefreshToken)
			assert.Nil(t, err)
			assert.NotEqual(t, refreshTokenUuid, uuid.Nil)

			if oldRefreshToken != nil {
				// make sure old login token doesn't exist
				oldTokenRes, err := refreshTokenRepo.GetByToken(context.Background(), oldRefreshToken.Token)
				assert.Empty(t, oldTokenRes)
				assert.Error(t, err)

				// make sure old refresh token and new refresh token dates are the same.
				newTokenRes, err := refreshTokenRepo.GetByToken(context.Background(), refreshTokenUuid)
				assert.NotEmpty(t, newTokenRes)
				assert.Nil(t, err)
				assert.Equal(t, oldRefreshToken.ValidUntil, newTokenRes.ValidUntil)
			}
		})
	}
}
