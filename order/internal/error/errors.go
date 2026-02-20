package error

import "errors"

var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrInvalidFormat   = errors.New("invalid format")
	ErrStatusPaid      = errors.New("order already paid")
	ErrStatusCancelled = errors.New("order already cancelled")
	ErrInvalidState    = errors.New("invalid state")
)
