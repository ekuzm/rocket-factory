package error

import "errors"

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidUUIDFormat = errors.New("invalid UUID format")
	ErrStatusPaid        = errors.New("order already paid")
	ErrStatusCancelled   = errors.New("order already cancelled")
)
