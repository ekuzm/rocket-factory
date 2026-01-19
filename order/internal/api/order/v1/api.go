package v1

import (
	"github.com/ekuzm/rocket-factory/order/internal/service"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

type api struct {
	service service.OrderService
	orderV1.UnimplementedHandler
}

func NewAPI(service service.OrderService) *api {
	return &api{
		service: service,
	}
}
