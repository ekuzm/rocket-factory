package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	api "github.com/ekuzm/rocket-factory/payment/internal/api/v1"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
)

var _ api.PaymentService = (*service)(nil)

type service struct{}

func New() *service {
	return &service{}
}

func (s *service) PayOrder(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	if paymentMethod == model.PaymentMethodUnknown {
		err := fmt.Errorf("pay order: %w", errs.ErrInvalid)
		logger.WithFields(logrus.Fields{
			"Order UUID":     orderUUID,
			"User UUID":      userUUID,
			"Payment Method": paymentMethod,
			"error":          err,
		}).Warn("Failed to pay order")

		return uuid.Nil, err
	}

	transactionUUID := uuid.New()

	logger.Debug("Payment was successfully, transaction uuid: ", transactionUUID)

	return transactionUUID, nil
}
