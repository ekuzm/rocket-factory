package v1

import (
	"context"

	"github.com/ekuzm/rocket-factory/order/internal/converter"
	"github.com/ekuzm/rocket-factory/order/internal/dto"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (*orderV1.Order, error) {
	input := dto.GetOrderInput{
		UUID: params.OrderUUID,
	}

	output, err := a.service.GetOrder(ctx, input)
	if err != nil {
		return nil, err // error handling in NewError(...) method
	}

	return converter.OrderToAPI(output.Order), nil
}
