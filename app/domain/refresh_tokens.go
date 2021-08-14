package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Refresh tokens for the JWTs
type RefreshToken struct {
	Id         uuid.UUID `json:"-"`
	Token      string    `json:"token"`
	ValidUntil time.Time `json:"-"`
}

type RefreshTokenRepository interface {
	GetByToken(ctx context.Context, token string) (RefreshToken, error)
	Store(ctx context.Context, token RefreshToken) (RefreshToken, error)
}

type RefreshTokenService interface {
	IsTokenValid(ctx context.Context, token string) (RefreshToken, error)
	GenerateRefreshToken(ctx context.Context, userId uuid.UUID) (RefreshToken, error)
}
