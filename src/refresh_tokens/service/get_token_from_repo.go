package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

// Gets the refresh token from the repo, based on an uuid token
func (s DefaultRefreshTokenService) GetTokenFromRepo(ctx context.Context, token uuid.UUID) (domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	refreshToken, err := s.TokenRepo.GetByToken(ctx, token)
	if err != nil {
		s.Logger.Printf("error fetching refresh token by value {%s}: %v\n", token.String(), err)
		return domain.RefreshToken{}, err
	}
	return refreshToken, nil
}
