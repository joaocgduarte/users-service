package tokens

import (
	"testing"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/domain/mocks"
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
		tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlTGFiZWwiOiJBZG1pbmlzdHJhdG9yIiwidXNlcm5hbWUiOiJhZG1pbiIsImV4cCI6MTYzMTEzODYyNywiaXNzIjoiYTNiMzI1MDQtYjc4Mi00NDU1LTg3M2YtZGM0OWQzMDdmMjgyIn0.FgFqj8ucoWLaskO1Akba22bJ5Kyaw6krToTr_k52UWc"
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
