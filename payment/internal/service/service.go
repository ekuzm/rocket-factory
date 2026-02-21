package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	api "github.com/ekuzm/rocket-factory/payment/internal/api/v1"
	errs "github.com/ekuzm/rocket-factory/payment/internal/error"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
)

var _ api.PaymentService = (*service)(nil)

type service struct{}

func New() *service {
	return &service{}
}

func (s *service) PayOrder(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	if paymentMethod == model.PaymentMethodUnknown {
		return uuid.Nil, fmt.Errorf("pay order: %w", errs.ErrInvalidPaymentMethod)
	}

	transactionUUID := uuid.New()

	log.Printf("Payment was successfully, transaction uuid: %s", transactionUUID)

	return transactionUUID, nil
}
