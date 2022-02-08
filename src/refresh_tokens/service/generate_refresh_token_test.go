package service

import (
	"context"
	"testing"
	"time"

	"github.com/plagioriginal/user-microservice/domain"
	"github.com/stretchr/testify/assert"
)

//@todo: all the other tests for this function. Maybe add a logger to the service

func TestGenerateRefreshToken_InvalidInput(t *testing.T) {
	res, err := New(nil, nil, time.Duration(5*time.Second)).GenerateRefreshToken(context.TODO(), nil)
	assert.Empty(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, domain.ErrBadParamInput)
}
