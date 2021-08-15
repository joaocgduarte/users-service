package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

type DefaultRefreshTokenService struct {
	TokenRepo      domain.RefreshTokenRepository
	UserRepo       domain.UserRepository
	ContextTimeout time.Duration
}

func New(tokenRepo domain.RefreshTokenRepository, userRepo domain.UserRepository, contextTimeout time.Duration) domain.RefreshTokenService {
	return DefaultRefreshTokenService{tokenRepo, userRepo, contextTimeout}
}

func (s DefaultRefreshTokenService) GetUserByToken(ctx context.Context, token domain.RefreshToken) (*domain.User, error) {
	_, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	user, err := s.UserRepo.GetUserByRefreshToken(ctx, token.Id)

	if err != nil {
		return &domain.User{}, err
	}

	return user, nil
}

// Gets the refresh token from the repo, based on an uuid token
func (s DefaultRefreshTokenService) GetTokenFromRepo(ctx context.Context, token uuid.UUID) (domain.RefreshToken, error) {
	_, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	refreshToken, err := s.TokenRepo.GetByToken(ctx, token)

	if err != nil {
		return domain.RefreshToken{}, err
	}

	return refreshToken, nil
}

// Deletes a token from the DB (basically logout).
func (s DefaultRefreshTokenService) DeleteToken(ctx context.Context, token domain.RefreshToken) error {
	_, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	err := s.TokenRepo.Delete(ctx, token.Id)

	if err != nil {
		return err
	}

	return nil
}

// Checks weather the refresh token is Valid.
func (s DefaultRefreshTokenService) IsTokenValid(token domain.RefreshToken) bool {
	return token.ValidUntil.After(time.Now())
}

// Generates the refresh tokens
func (s DefaultRefreshTokenService) GenerateRefreshToken(ctx context.Context, user *domain.User) (domain.RefreshToken, error) {
	_, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	if user.RefreshTokenId.Valid {
		err := s.TokenRepo.Delete(ctx, user.RefreshTokenId.UUID)

		if err != nil {
			return domain.RefreshToken{}, err
		}
	}

	newToken := uuid.New()
	validUntil := time.Now().Add(time.Hour * 24 * 7)

	refreshToken := domain.RefreshToken{
		Token:      newToken,
		ValidUntil: validUntil,
	}

	refreshToken, err := s.TokenRepo.Store(ctx, refreshToken)

	if err != nil {
		return domain.RefreshToken{}, err
	}

	err = s.UserRepo.SaveRefreshToken(ctx, user, refreshToken)

	if err != nil {
		return domain.RefreshToken{}, err
	}

	return refreshToken, nil
}
