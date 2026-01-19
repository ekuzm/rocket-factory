package repository

import (
	"context"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	repoModel "github.com/ekuzm/rocket-factory/order/internal/repository/model"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order repoModel.Order) error
	PayOrder(ctx context.Context, order repoModel.Order) error
	GetOrder(ctx context.Context, uuid string) (model.Order, error)
	CancelOrder(ctx context.Context, order repoModel.Order) error
}
