package client

import (
	"context"

	"github.com/ekuzm/rocket-factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, payment model.Payment) (string, error)
}
