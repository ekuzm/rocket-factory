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
	if _, err := uuid.Parse(orderUUID); err != nil {
		return Payment{}, fmt.Errorf("order uuid: %w", ErrInvalidFormat)
	}
	if _, err := uuid.Parse(userUUID); err != nil {
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

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSPB           PaymentMethod = "SPB"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)
