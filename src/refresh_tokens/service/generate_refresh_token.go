package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

// Generates the refresh tokens
func (s DefaultRefreshTokenService) GenerateRefreshToken(ctx context.Context, user *domain.User) (domain.RefreshToken, error) {
	if user == nil {
		return domain.RefreshToken{}, domain.ErrBadParamInput
	}

	ctx, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	validUntil := time.Now().Add(time.Hour * 24 * 7)

	if user.RefreshTokenId.Valid {
		oldToken, err := s.TokenRepo.GetByUUID(ctx, user.RefreshTokenId.UUID)
		if err != nil {
			return domain.RefreshToken{}, err
		}

		validUntil = oldToken.ValidUntil
		err = s.TokenRepo.Delete(ctx, oldToken.Id)
		if err != nil {
			return domain.RefreshToken{}, err
		}
	}

	newToken := uuid.New()
	refreshToken := domain.RefreshToken{
		Token:      newToken,
		ValidUntil: validUntil,
	}

	refreshToken, err := s.TokenRepo.Store(ctx, refreshToken)
	if err != nil {
		return domain.RefreshToken{}, err
	}

	if err = s.UserRepo.SaveRefreshToken(ctx, user, refreshToken); err != nil {
		if deleteErr := s.TokenRepo.Delete(ctx, refreshToken.Id); deleteErr != nil {
			return domain.RefreshToken{}, deleteErr
		}
		return domain.RefreshToken{}, err
	}
	return refreshToken, nil
}
