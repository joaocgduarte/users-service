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

// Checks weather the refresh token is Valid. If it is not, it will return an error.
// Returns the token object from DB if it is valid.
func (s DefaultRefreshTokenService) IsTokenValid(ctx context.Context, token uuid.UUID) (domain.RefreshToken, error) {
	_, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	refreshToken, err := s.TokenRepo.GetByToken(ctx, token)

	if err != nil {
		return domain.RefreshToken{}, err
	}

	if !refreshToken.ValidUntil.After(time.Now()) {
		s.TokenRepo.Delete(ctx, refreshToken.Id)
		return domain.RefreshToken{}, domain.ErrNotFound
	}

	return refreshToken, nil
}

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
