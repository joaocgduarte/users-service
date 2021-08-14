package tokens

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

// Our custom claimes for the JWT Token
type ClaimsWithRole struct {
	UserRole string `json:"role"`
	jwt.StandardClaims
}

// Object used to manage token/auth operations
type TokenManager struct {
	JWTSecret string
}

// Instantiates a new Token Manager
func New(jwtSecret string) TokenManager {
	return TokenManager{JWTSecret: jwtSecret}
}

// Generates a new JWT token for a given user.
func (t TokenManager) GenerateJWT(user domain.User) (string, error) {
	claims := ClaimsWithRole{
		user.Role.RoleSlug,
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
		return claims.UserRole, nil
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
