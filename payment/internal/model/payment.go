package model

import (
	"fmt"

	"github.com/google/uuid"
)

type Payment struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethod
}

func NewPayment(orderUUID, userUUID string, paymentMethod PaymentMethod) (Payment, error) {
	if err := uuid.Validate(orderUUID); err != nil {
		return Payment{}, fmt.Errorf("order uuid: %w", ErrInvalidFormat)
	}
	if err := uuid.Validate(userUUID); err != nil {
		return Payment{}, fmt.Errorf("user uuid: %w", ErrInvalidFormat)
	}
	if paymentMethod == PaymentMethodUnknown {
		return Payment{}, fmt.Errorf("payment method: %w", ErrInvalidFormat)
	}

	return Payment{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PaymentMethod: paymentMethod,
	}, nil
}

type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSPB
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)
