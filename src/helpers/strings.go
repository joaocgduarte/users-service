package helpers

import "context"

func StringFromContext(ctx context.Context, key string) string {
	val := ctx.Value(key)
	if val == nil {
		return ""
	}
	res, ok := val.(string)
	if !ok {
		return ""
	}
	return res
}
