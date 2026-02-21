package model

import (
	"errors"
)

var (
	ErrInvalidUUIDFormat error = errors.New("invalid UUID format")
	ErrPartNotFound      error = errors.New("part not found")
)
