package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Refresh tokens for the JWTs
type RefreshToken struct {
	Id         uuid.UUID `json:"-"`
	Token      uuid.UUID `json:"token"`
	ValidUntil time.Time `json:"-"`
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
