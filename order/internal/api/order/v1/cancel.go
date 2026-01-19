package v1

import (
	"context"
	"net/http"

	"github.com/ekuzm/rocket-factory/order/internal/dto"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (*orderV1.NoContent, error) {
	input := dto.CancelOrderInput{
		UUID: params.OrderUUID,
	}

	err := a.service.CancelOrder(ctx, input)
	if err != nil {
		return nil, err // error handling in NewError(...) method
	}

	return &orderV1.NoContent{Code: http.StatusNoContent, Message: "order successfully cancelled"}, nil
}
