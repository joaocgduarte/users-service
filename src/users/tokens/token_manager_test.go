package tokens

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/domain/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Define the suite
type TokenManagerTestSuite struct {
	suite.Suite
	jwtSecret           string
	refreshTokenService *mocks.RefreshTokenService
	roleRepo            *mocks.RoleRepository
	validMockUser       *domain.User
	invalidMockUsers    []*domain.User
}

// Runs before all the tests
func (ts *TokenManagerTestSuite) SetupSuite() {
	ts.jwtSecret = "jwt-generator-secret"
	userId, roleId := uuid.New(), uuid.New()

	ts.validMockUser = &domain.User{
		ID:       userId,
		Username: "testuser",
		RoleId:   roleId,
		Role: &domain.Role{
			ID:        roleId,
			RoleSlug:  domain.DEFAULT_ROLE_ADMIN.RoleSlug,
			RoleLabel: domain.DEFAULT_ROLE_ADMIN.RoleLabel,
		},
	}

	ts.invalidMockUsers = []*domain.User{
		{
			ID:       userId,
			Username: "",
			RoleId:   roleId,
			Role: &domain.Role{
				ID:        roleId,
				RoleSlug:  domain.DEFAULT_ROLE_ADMIN.RoleSlug,
				RoleLabel: domain.DEFAULT_ROLE_ADMIN.RoleLabel,
			},
		},
		{
			ID:       uuid.Nil,
			Username: "admin",
			RoleId:   roleId,
			Role: &domain.Role{
				ID:        roleId,
				RoleSlug:  domain.DEFAULT_ROLE_ADMIN.RoleSlug,
				RoleLabel: domain.DEFAULT_ROLE_ADMIN.RoleLabel,
			},
		},
		{
			ID:       userId,
			Username: "admin",
			RoleId:   uuid.Nil,
			Role: &domain.Role{
				ID:        uuid.Nil,
				RoleSlug:  domain.DEFAULT_ROLE_ADMIN.RoleSlug,
				RoleLabel: domain.DEFAULT_ROLE_ADMIN.RoleLabel,
			},
		},
		{
			ID:       userId,
			Username: "admin",
			RoleId:   roleId,
			Role: &domain.Role{
				ID:        roleId,
				RoleSlug:  "",
				RoleLabel: domain.DEFAULT_ROLE_ADMIN.RoleLabel,
			},
		},
		{
			ID:       userId,
			Username: "admin",
			RoleId:   roleId,
			Role: &domain.Role{
				ID:        roleId,
				RoleSlug:  domain.DEFAULT_ROLE_ADMIN.RoleSlug,
				RoleLabel: "",
			},
		},
	}
}

// Runs before each test.
func (ts *TokenManagerTestSuite) SetupTest() {
	ts.refreshTokenService = new(mocks.RefreshTokenService)
	ts.roleRepo = new(mocks.RoleRepository)
}

// Tests the generate refresh token function.
func (ts *TokenManagerTestSuite) TestGenerateRefreshToken() {
	ts.Run("invalid user to generate token", func() {
		ts.refreshTokenService.
			On("GenerateRefreshToken", mock.Anything, ts.validMockUser).
			Return(domain.RefreshToken{}, errors.New("any error")).
			Once()

		tm := NewTokenManager(ts.jwtSecret, ts.refreshTokenService, ts.roleRepo)

		token, err := tm.GenerateRefreshToken(context.TODO(), ts.validMockUser)

		ts.Equal(token, uuid.Nil)
		ts.Error(err)
		ts.refreshTokenService.AssertExpectations(ts.T())
	})

	ts.Run("valid user to generate token", func() {
		ts.refreshTokenService.
			On("GenerateRefreshToken", mock.Anything, ts.validMockUser).
			Return(domain.RefreshToken{
				Id:         uuid.New(),
				Token:      uuid.New(),
				ValidUntil: time.Now().Add(time.Hour * 24 * 7),
			}, nil).
			Once()

		tm := NewTokenManager(ts.jwtSecret, ts.refreshTokenService, ts.roleRepo)

		token, err := tm.GenerateRefreshToken(context.TODO(), ts.validMockUser)

		ts.NotEqual(token, uuid.Nil)
		ts.NoError(err)
		ts.refreshTokenService.AssertExpectations(ts.T())
	})
}

// Tests the ID and role getters for a JWT.
func (ts *TokenManagerTestSuite) TestGettersFromJWT() {
	tm := NewTokenManager(ts.jwtSecret, ts.refreshTokenService, ts.roleRepo)
	tokenString, _ := tm.GenerateJWT(ts.validMockUser)
	token, _ := tm.ParseJWT(tokenString)

	ts.Run("User ID getter works", func() {
		id, err := tm.GetUserIDFromToken(token)
		ts.NoError(err)
		ts.Equal(id, ts.validMockUser.ID)
	})

	ts.Run("user role getter works", func() {
		role, err := tm.GetUserRoleFromToken(token)
		ts.NoError(err)
		ts.Equal(role, ts.validMockUser.Role.RoleSlug)
	})

	token.Valid = false

	ts.Run("User ID getter works but token is invalid", func() {
		id, err := tm.GetUserIDFromToken(token)
		ts.Error(err)
		ts.NotEqual(id, ts.validMockUser.ID)
	})

	ts.Run("user role getter works", func() {
		role, err := tm.GetUserRoleFromToken(token)
		ts.Error(err)
		ts.NotEqual(role, ts.validMockUser.Role.RoleSlug)
	})
}

// Tests the parsing of the JWT tokens
func (ts *TokenManagerTestSuite) TestParseJWT() {
	tm := NewTokenManager(ts.jwtSecret, ts.refreshTokenService, ts.roleRepo)

	ts.Run("valid structure jwt", func() {
		tokenString, _ := tm.GenerateJWT(ts.validMockUser)
		token, err := tm.ParseJWT(tokenString)

		ts.NoError(err)
		ts.True(token.Valid)

		validClaims := token.Claims.Valid()
		ts.NoError(validClaims)
	})

	ts.Run("invalid signature jwt", func() {
		tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlTGFiZWxzIjoiQWRtaW5pc3RyYXRvciIsInVzZXJuYW1lIjoiYWRtaW4iLCJleHAiOjE2MzExMzg2MjcsImlzcyI6ImEzYjMyNTA0LWI3ODItNDQ1NS04NzNmLWRjNDlkMzA3ZjI4MiJ9.pgMS23a7YjDuJIBczKdSarvPlQXPxDOnyuUJrHeXccI"
		token, err := tm.ParseJWT(tokenString)

		ts.Error(err)
		ts.False(token.Valid)
	})
}

// Tests the generation of the jwts
func (ts *TokenManagerTestSuite) TestGenerateJWT() {
	tm := NewTokenManager(ts.jwtSecret, ts.refreshTokenService, ts.roleRepo)

	ts.Run("user is valid and generates proper token", func() {
		token, err := tm.GenerateJWT(ts.validMockUser)

		ts.NoError(err)
		ts.NotEmpty(token)
		ts.NotEqual(token, "")
	})

	ts.Run("user is invalid and doesn't generate token", func() {
		for _, invalidUser := range ts.invalidMockUsers {
			token, err := tm.GenerateJWT(invalidUser)

			ts.Empty(token)
			ts.Equal("", token)
			ts.Error(err)
		}
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TokenManagerTestSuite))
}
