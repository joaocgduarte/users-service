package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Refresh tokens for the JWTs
type RefreshToken struct {
	Id         uuid.UUID `json:"-"`
	Token      uuid.UUID `json:"token"`
	ValidUntil time.Time `json:"-"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AccessTokenHandler interface {
	ParseJWT(tokenString string) (*jwt.Token, error)
	IsJWTokenValid(token *jwt.Token) bool
	GetUserRoleFromToken(token *jwt.Token) (string, error)
	GenerateTokens(ctx context.Context, user *User) (TokenResponse, error)
	RefreshAllTokens(ctx context.Context, askedRefreshToken uuid.UUID) (TokenResponse, error)
	GetUserIDFromToken(token *jwt.Token) (uuid.UUID, error)
	DeleteRefreshToken(ctx context.Context, refreshToken string) bool
}

type RefreshTokenRepository interface {
	GetByUUID(ctx context.Context, id uuid.UUID) (RefreshToken, error)
	GetByToken(ctx context.Context, token uuid.UUID) (RefreshToken, error)
	Store(ctx context.Context, token RefreshToken) (RefreshToken, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type RefreshTokenService interface {
	GetUserByToken(ctx context.Context, token RefreshToken) (*User, error)
	DeleteToken(ctx context.Context, token RefreshToken) error
	GetTokenFromRepo(ctx context.Context, token uuid.UUID) (RefreshToken, error)
	IsTokenValid(token RefreshToken) bool
	GenerateRefreshToken(ctx context.Context, user *User) (RefreshToken, error)
}
