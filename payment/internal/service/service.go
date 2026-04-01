package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	api "github.com/ekuzm/rocket-factory/payment/internal/api/v1"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
)

var _ api.PaymentService = (*service)(nil)

type service struct{}

func New() *service {
	return &service{}
}

func (s *service) PayOrder(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	if paymentMethod == model.PaymentMethodUnknown {
		err := fmt.Errorf("pay order: %w", errs.ErrInvalid)
		slog.Warn(
			"payment rejected",
			"orderUUID", orderUUID,
			"userUUID", userUUID,
			"paymentMethod", paymentMethod,
			"error", err,
		)

		return uuid.Nil, err
	}

	transactionUUID := uuid.New()

	slog.Debug("payment succeeded", "transactionUUID", transactionUUID)

	return transactionUUID, nil
}
