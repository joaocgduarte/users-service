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
	GetByToken(ctx context.Context, token uuid.UUID) (RefreshToken, error)
	Store(ctx context.Context, token RefreshToken) (RefreshToken, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type RefreshTokenService interface {
	IsTokenValid(ctx context.Context, token uuid.UUID) (RefreshToken, error)
	GenerateRefreshToken(ctx context.Context, user *User) (RefreshToken, error)
}
