package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

type adapter struct {
	grpcClient paymentV1.PaymentServiceClient
}

func New(grpcClient paymentV1.PaymentServiceClient) *adapter {
	return &adapter{
		grpcClient: grpcClient,
	}
}

func (a *adapter) PayOrder(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	req := &paymentV1.PayOrderRequest{
		Uuid:          orderUUID.String(),
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentMethodToGRPC[paymentMethod],
	}

	resp, err := a.grpcClient.PayOrder(ctx, req)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Order UUID":     orderUUID,
			"User UUID":      userUUID,
			"Payment Method": paymentMethod,
			"error":          err,
		}).Error("Failed to pay order in payment service")

		return uuid.Nil, fmt.Errorf("pay order: %w", err)
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Order UUID":       orderUUID,
			"User UUID":        userUUID,
			"Payment Method":   paymentMethod,
			"Transaction UUID": resp.TransactionUuid,
			"error":            err,
		}).Error("Failed to parse transaction UUID")

		return uuid.Nil, fmt.Errorf("transaction UUID: %w", errs.ErrInvalid)
	}

	return transactionUUID, nil
}
