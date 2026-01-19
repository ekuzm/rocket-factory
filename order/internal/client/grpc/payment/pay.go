package payment

import (
	"context"

	clientConverter "github.com/ekuzm/rocket-factory/order/internal/client/converter"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, payment model.Payment) (string, error) {
	req := &paymentV1.PayOrderRequest{
		Uuid:          payment.OrderUUID,
		UserUuid:      payment.UserUUID,
		PaymentMethod: clientConverter.PaymentMethodToProto(payment.PaymentMethod),
	}

	resp, err := c.generatedClient.PayOrder(ctx, req)
	if err != nil {
		return "", nil
	}

	return resp.TransactionUuid, nil
}
