package service

import (
	"context"

	"github.com/plagioriginal/user-microservice/domain"
)

// Deletes a token from the DB (basically logout).
func (s DefaultRefreshTokenService) DeleteToken(ctx context.Context, token domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	return s.TokenRepo.Delete(ctx, token.Id)
}
