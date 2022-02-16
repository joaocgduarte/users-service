package service

import (
	"time"

	"github.com/plagioriginal/user-microservice/domain"
)

type DefaultUserService struct {
	UserRepo       domain.UserRepository
	RoleRepo       domain.RoleRepository
	ContextTimeout time.Duration
}

// Constructor
func New(
	userRepo domain.UserRepository,
	roleRepo domain.RoleRepository,
	contextTimeout time.Duration,
) domain.UserService {
	return DefaultUserService{
		userRepo,
		roleRepo,
		contextTimeout,
	}
}

// Used for tests
func newService(userRepo domain.UserRepository, roleRepo domain.RoleRepository) domain.UserService {
	return New(
		userRepo,
		roleRepo,
		time.Duration(5*time.Second),
	)
}
