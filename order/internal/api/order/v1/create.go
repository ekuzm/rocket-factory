package v1

import (
	"context"

	"github.com/ekuzm/rocket-factory/order/internal/dto"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (*orderV1.CreateOrderResponse, error) {
	input := dto.CreateOrderInput{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUuids,
	}

	output, err := a.service.CreateOrder(ctx, input)
	if err != nil {
		return nil, err // error handling in NewError(...) method
	}

	return &orderV1.CreateOrderResponse{UUID: output.OrderUUID, TotalPrice: output.TotalPrice}, nil
}
