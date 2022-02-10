package service

import (
	"testing"
	"time"

	"github.com/plagioriginal/user-microservice/domain"
	"github.com/stretchr/testify/assert"
)

func TestIsTokenValid_InvalidToken(t *testing.T) {
	isValid := newService(nil, nil).IsTokenValid(domain.RefreshToken{
		ValidUntil: time.Now().Add(time.Duration(-5) * time.Hour),
	})
	assert.False(t, isValid)
}

func TestIsTokenValid_ValidToken(t *testing.T) {
	isValid := newService(nil, nil).IsTokenValid(domain.RefreshToken{
		ValidUntil: time.Now().Add(time.Duration(5) * time.Hour),
	})
	assert.True(t, isValid)
}
