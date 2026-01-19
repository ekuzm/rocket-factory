package v1

import (
	"context"

	"github.com/ekuzm/rocket-factory/order/internal/dto"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (*orderV1.PayOrderResponse, error) {
	input := dto.PayOrderInput{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
	}

	output, err := a.service.PayOrder(ctx, input)
	if err != nil {
		return nil, err // error handling in NewError(...) method
	}

	return &orderV1.PayOrderResponse{TransactionUUID: output.TransactionUUID}, nil
}
