package model

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidFormat = errors.New("invalid format")
	ErrConflict      = errors.New("conflict")
)
