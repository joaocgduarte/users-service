package domain

import "errors"

var (
	ErrBadParamInput = errors.New("invalid parameter")
	ErrNotFound      = errors.New("resource not found")
)
