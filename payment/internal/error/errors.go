package error

import "errors"

var (
	ErrInvalidUUID          = errors.New("invalid uuid format")
	ErrInvalidPaymentMethod = errors.New("payment method is an unknown")
)
