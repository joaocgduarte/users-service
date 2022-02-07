package service

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/plagioriginal/user-microservice/domain"
)

type DefaultUserService struct {
	Logger         *log.Logger
	UserRepo       domain.UserRepository
	RoleRepo       domain.RoleRepository
	ContextTimeout time.Duration
}

// Constructor
func New(
	logger *log.Logger,
	userRepo domain.UserRepository,
	roleRepo domain.RoleRepository,
	contextTimeout time.Duration,
) domain.UserService {
	return DefaultUserService{
		logger,
		userRepo,
		roleRepo,
		contextTimeout,
	}
}

// Used for tests
func newService(userRepo domain.UserRepository, roleRepo domain.RoleRepository) domain.UserService {
	return New(
		log.New(ioutil.Discard, "tests: ", log.Flags()),
		userRepo,
		roleRepo,
		time.Duration(5*time.Second),
	)
}
