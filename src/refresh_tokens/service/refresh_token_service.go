package service

import (
	"time"

	"github.com/plagioriginal/user-microservice/domain"
)

type DefaultRefreshTokenService struct {
	TokenRepo      domain.RefreshTokenRepository
	UserRepo       domain.UserRepository
	ContextTimeout time.Duration
}

// New service Instantiation
func New(
	tokenRepo domain.RefreshTokenRepository,
	userRepo domain.UserRepository,
	contextTimeout time.Duration,
) domain.RefreshTokenService {
	return DefaultRefreshTokenService{tokenRepo, userRepo, contextTimeout}
}
