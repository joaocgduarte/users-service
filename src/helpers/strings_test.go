package helpers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromContext_NilContextReturnEmptyString(t *testing.T) {
	res := StringFromContext(context.TODO(), "cenas")
	assert.Empty(t, res)
}

func TestFromContext_Success(t *testing.T) {
	ctx := context.WithValue(context.TODO(), "cenas", "valor")
	res := StringFromContext(ctx, "cenas")
	assert.Equal(t, "valor", res)
}

func TestFromContext_DifferentStructure(t *testing.T) {
	ctx := context.WithValue(context.TODO(), "cenas", 123)
	res := StringFromContext(ctx, "cenas")
	assert.Equal(t, "", res)
}
