package model

import (
	"errors"
)

var (
	ErrInvalidFormat error = errors.New("invalid format")
	ErrNotFound      error = errors.New("not found")
)
