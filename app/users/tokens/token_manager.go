package tokens

import (
	"context"
	"errors"
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
}

// Instantiates a new Token Manager
func NewTokenManager(jwtSecret string, refreshTokenService domain.RefreshTokenService) TokenManager {
	return TokenManager{JWTSecret: jwtSecret, RefreshTokenService: refreshTokenService}
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

	return uuid.Nil, errors.New("invalid token")
}

// Gets the user role from a jwt token
func (t TokenManager) GetUserRoleFromToken(token *jwt.Token) (string, error) {
	if claims, ok := token.Claims.(*ClaimsWithRole); ok && t.IsJWTokenValid(token) {
		return claims.UserRoleSlug, nil
	}

	return "", errors.New("invalid token")
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
