package v1

import (
	"context"

	"github.com/ekuzm/rocket-factory/payment/internal/dto"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	input := dto.PayOrderInput{
		OrderUUID:     req.GetUuid(),
		UserUUID:      req.GetUserUuid(),
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
	}

	output, err := a.service.PayOrder(ctx, input)
	if err != nil {
		return nil, err // error handling in mapping errors interceptor
	}

	return &paymentV1.PayOrderResponse{TransactionUuid: output.TransactionUUID}, nil
}
