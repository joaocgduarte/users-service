package service

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/plagioriginal/user-microservice/domain"
)

type DefaultRefreshTokenService struct {
	Logger         *log.Logger
	TokenRepo      domain.RefreshTokenRepository
	UserRepo       domain.UserRepository
	ContextTimeout time.Duration
}

// New service Instantiation
func New(
	logger *log.Logger,
	tokenRepo domain.RefreshTokenRepository,
	userRepo domain.UserRepository,
	contextTimeout time.Duration,
) domain.RefreshTokenService {
	return DefaultRefreshTokenService{logger, tokenRepo, userRepo, contextTimeout}
}

// Instantiation for tests
func newService(
	tokenRepo domain.RefreshTokenRepository,
	userRepo domain.UserRepository,
) domain.RefreshTokenService {
	return DefaultRefreshTokenService{
		log.New(ioutil.Discard, "tests: ", log.Flags()),
		tokenRepo,
		userRepo,
		time.Duration(5 * time.Second),
	}
}
