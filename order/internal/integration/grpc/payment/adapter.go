package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/ekuzm/rocket-factory/order/internal/error"
	"github.com/ekuzm/rocket-factory/order/internal/model"
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
		return uuid.Nil, fmt.Errorf("pay order: %w", err)
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("transaction UUID: %w", errs.ErrInvalidUUIDFormat)
	}

	return transactionUUID, nil
}
