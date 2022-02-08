package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID        uuid.UUID `json:"id"`
	RoleSlug  string    `json:"roleSlug"`
	RoleLabel string    `json:"roleLabel"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"-"`
}

type RoleRepository interface {
	Fetch(ctx context.Context) ([]Role, error)
	Store(ctx context.Context, role Role) (Role, error)
	GetBySlug(ctx context.Context, slug string) (Role, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (Role, error)
}

type RoleService interface {
	Store(ctx context.Context, role Role) (Role, error)
	GetBySlug(ctx context.Context, slug string) (Role, error)
}

var (
	DEFAULT_ROLE_ADMIN = Role{
		RoleSlug:  "admin",
		RoleLabel: "Administrator",
	}
	DEFAULT_ROLE_USER = Role{
		RoleSlug:  "user",
		RoleLabel: "User",
	}
)
