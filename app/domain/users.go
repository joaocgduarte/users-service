package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID     `json:"id"`
	FirstName      string        `json:"firstName"`
	LastName       string        `json:"lastName"`
	Username       string        `json:"username"`
	Password       string        `json:"-"`
	RoleId         uuid.UUID     `json:"-"`
	Role           *Role         `json:"role,omitempty"`
	RefreshTokenId uuid.NullUUID `json:"-"`
	RefreshToken   string        `json:"-"`
	CreatedAt      time.Time     `json:"-"`
	UpdatedAt      time.Time     `json:"-"`
	DeletedAt      time.Time     `json:"-"`
}

type UserRepository interface {
	Store(ctx context.Context, user User) (*User, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	SaveRefreshToken(ctx context.Context, user *User, token RefreshToken) error
}

type UserService interface {
	Store(ctx context.Context, username string, password string, roleSlug string) (*User, error)
	GetUserByLogin(ctx context.Context, username string, password string) (*User, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*User, error)
}
