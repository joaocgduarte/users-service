package service

import (
	"time"

	"github.com/plagioriginal/user-microservice/domain"
)

const (
	ProductionBcryptCost int = 14
	TestingBcryptCost    int = 2
)

type DefaultUserService struct {
	UserRepo          domain.UserRepository
	RoleRepo          domain.RoleRepository
	ContextTimeout    time.Duration
	BcryptHashingCost int
}

// Constructor
func New(
	userRepo domain.UserRepository,
	roleRepo domain.RoleRepository,
	contextTimeout time.Duration,
	bcryptHashingCost int,
) domain.UserService {
	return DefaultUserService{
		userRepo,
		roleRepo,
		contextTimeout,
		bcryptHashingCost,
	}
}

// Used for tests
func newService(userRepo domain.UserRepository, roleRepo domain.RoleRepository) domain.UserService {
	return New(
		userRepo,
		roleRepo,
		time.Duration(2*time.Second),
		2,
	)
}
