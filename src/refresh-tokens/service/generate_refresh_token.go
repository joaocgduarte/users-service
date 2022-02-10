package service

import (
	"context"
	"errors"
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
			s.Logger.Printf("error fetching token by id even though it's in user's {%s} reference: %v\n", user.ID.String(), err)
			return domain.RefreshToken{}, err
		}
		if !s.IsTokenValid(oldToken) {
			s.Logger.Printf("user {%s} refresh token is invalid: %v\n", user.ID.String(), oldToken)
			return domain.RefreshToken{}, errors.New("invalid refresh token")
		}

		validUntil = oldToken.ValidUntil
		err = s.TokenRepo.Delete(ctx, oldToken.Id)
		if err != nil {
			s.Logger.Printf("error deleting user {%s} refresh token: %v\n", user.ID.String(), err)
			return domain.RefreshToken{}, err
		}
	}

	newToken := uuid.New()
	refreshTokenIn := domain.RefreshToken{
		Token:      newToken,
		ValidUntil: validUntil,
	}

	refreshToken, err := s.TokenRepo.Store(ctx, refreshTokenIn)
	if err != nil {
		s.Logger.Printf("error storing new refresh token {%v} refresh token: %v\n", refreshTokenIn, err)
		return domain.RefreshToken{}, err
	}

	if err = s.UserRepo.SaveRefreshToken(ctx, user, refreshToken); err != nil {
		s.Logger.Printf("error saving refresh token {%v} reference in users {%s} table: %v\n", refreshTokenIn, user.ID.String(), err)
		if deleteErr := s.TokenRepo.Delete(ctx, refreshToken.Id); deleteErr != nil {
			s.Logger.Printf("error deleting refresh token {%v}: %v\n", refreshTokenIn, err)
			return domain.RefreshToken{}, deleteErr
		}
		return domain.RefreshToken{}, err
	}
	return refreshToken, nil
}
