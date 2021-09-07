package domain

import "errors"

var (
	ErrNotAllowed    = errors.New("not allowed")
	ErrBadParamInput = errors.New("invalid parameter")
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidToken  = errors.New("invalid token")
)
