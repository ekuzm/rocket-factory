package error

import "errors"

var (
	ErrInvalidUUID          = errors.New("invalid UUID format")
	ErrInvalidPaymentMethod = errors.New("payment method is an unknown")
)
