package tokens

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

// Our custom claimes for the JWT Token
type ClaimsWithRole struct {
	UserRoleSlug  string `json:"roleSlug"`
	UserRoleLabel string `json:"roleLabel"`
	Username      string `json:"username"`
	jwt.StandardClaims
}

// Object used to manage token/auth operations
type TokenManager struct {
	JWTSecret           string
	RefreshTokenService domain.RefreshTokenService
	RoleRepo            domain.RoleRepository
}

// Instantiates a new Token Manager
func NewTokenManager(jwtSecret string, refreshTokenService domain.RefreshTokenService, rolesRepo domain.RoleRepository) TokenManager {
	return TokenManager{
		JWTSecret:           jwtSecret,
		RefreshTokenService: refreshTokenService,
		RoleRepo:            rolesRepo,
	}
}

// Refreshes all the tokens, based on an old refresh token.
// If old token is invalid (out of date) then it will delete it from the DB and return error.
// If it is valid, it will generate new Access token and Refresh token to be used on next request.
func (t TokenManager) RefreshAllTokens(ctx context.Context, askedRefreshToken uuid.UUID) (domain.TokenResponse, error) {
	oldRefreshToken, err := t.RefreshTokenService.GetTokenFromRepo(ctx, askedRefreshToken)

	if err != nil {
		return domain.TokenResponse{}, err
	}

	isValid := t.RefreshTokenService.IsTokenValid(oldRefreshToken)

	if !isValid {
		err := t.RefreshTokenService.DeleteToken(ctx, oldRefreshToken)

		if err != nil {
			return domain.TokenResponse{}, err
		}

		return domain.TokenResponse{}, domain.ErrInvalidToken
	}

	user, err := t.RefreshTokenService.GetUserByToken(ctx, oldRefreshToken)

	if err != nil {
		return domain.TokenResponse{}, domain.ErrInvalidToken
	}

	userRole, err := t.RoleRepo.GetByUUID(ctx, user.RoleId)

	if err != nil {
		return domain.TokenResponse{}, err
	}

	user.Role = &userRole

	return t.GenerateTokens(ctx, user)
}

// Gets all the tokens as token response.
func (t TokenManager) GenerateTokens(ctx context.Context, user *domain.User) (domain.TokenResponse, error) {
	jwtToken, err := t.GenerateJWT(user)

	if err != nil {
		return domain.TokenResponse{}, err
	}

	refreshToken, err := t.GenerateRefreshToken(ctx, user)

	if err != nil {
		return domain.TokenResponse{}, err
	}

	return domain.TokenResponse{
		AccessToken:  jwtToken,
		RefreshToken: refreshToken.String(),
	}, nil
}

// Gets the regresh token
func (t TokenManager) GenerateRefreshToken(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	refreshToken, err := t.RefreshTokenService.GenerateRefreshToken(ctx, user)

	if err != nil {
		return uuid.Nil, err
	}

	return refreshToken.Token, nil
}

// Generates a new JWT token for a given user.
func (t TokenManager) GenerateJWT(user *domain.User) (string, error) {
	claims := ClaimsWithRole{
		user.Role.RoleSlug,
		user.Role.RoleLabel,
		user.Username,
		jwt.StandardClaims{
			Issuer:    user.ID.String(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := jwtToken.SignedString([]byte(t.JWTSecret))

	if err != nil {
		return "", err
	}

	return signedString, nil
}

// Gets the user ID from token
func (t TokenManager) GetUserIDFromToken(token *jwt.Token) (uuid.UUID, error) {
	if claims, ok := token.Claims.(*ClaimsWithRole); ok && t.IsJWTokenValid(token) {
		return uuid.Parse(claims.StandardClaims.Issuer)
	}

	return uuid.Nil, domain.ErrInvalidToken
}

// Gets the user role from a jwt token
func (t TokenManager) GetUserRoleFromToken(token *jwt.Token) (string, error) {
	if claims, ok := token.Claims.(*ClaimsWithRole); ok && t.IsJWTokenValid(token) {
		return claims.UserRoleSlug, nil
	}

	return "", domain.ErrInvalidToken
}

// Checks if a token is valid
func (t TokenManager) IsJWTokenValid(token *jwt.Token) bool {
	return token.Valid
}

// Parses a JWT Token string to an object.
func (t TokenManager) ParseJWT(tokenString string) (*jwt.Token, error) {
	key := []byte(t.JWTSecret)

	token, err := jwt.ParseWithClaims(tokenString, &ClaimsWithRole{}, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
