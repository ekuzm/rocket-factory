package service

import (
	"context"

	"github.com/ekuzm/rocket-factory/order/internal/dto"
)

type OrderService interface {
	CreateOrder(ctx context.Context, input dto.CreateOrderInput) (dto.CreateOrderOutput, error)
	GetOrder(ctx context.Context, input dto.GetOrderInput) (dto.GetOrderOutput, error)
	CancelOrder(ctx context.Context, input dto.CancelOrderInput) error
	PayOrder(ctx context.Context, input dto.PayOrderInput) (dto.PayOrderOutput, error)
}
