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
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	DeletedAt      time.Time     `json:"-"`
}

type StoreUserRequest struct {
	Username string
	Password string
	RoleSlug string
}

type GetUserRequest struct {
	Username string
	Password string
}

type UserRepository interface {
	Store(ctx context.Context, user User) (*User, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	SaveRefreshToken(ctx context.Context, user *User, token RefreshToken) error
	GetUserByRefreshToken(ctx context.Context, id uuid.UUID) (*User, error)
}

type UserService interface {
	Store(ctx context.Context, request StoreUserRequest) (*User, error)
	GetUserByLogin(ctx context.Context, request GetUserRequest) (*User, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*User, error)
}
