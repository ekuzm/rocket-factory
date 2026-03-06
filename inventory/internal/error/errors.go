package model

import (
	"errors"
)

var (
	ErrInvalidFilter     error = errors.New("invalid filter")
	ErrInvalidUUIDFormat error = errors.New("invalid UUID format")
	ErrPartNotFound      error = errors.New("part not found")
	ErrPartsNotFound     error = errors.New("some parts not found")
)
