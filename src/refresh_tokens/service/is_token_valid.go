package service

import (
	"time"

	"github.com/plagioriginal/user-microservice/domain"
)

// Checks weather the refresh token is Valid.
func (s DefaultRefreshTokenService) IsTokenValid(token domain.RefreshToken) bool {
	return token.ValidUntil.After(time.Now())
}
