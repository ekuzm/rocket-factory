package memory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	errs "github.com/ekuzm/rocket-factory/order/internal/error"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	"github.com/ekuzm/rocket-factory/order/internal/service"
)

var _ service.OrderRepository = (*repository)(nil)

type repository struct {
	orders map[uuid.UUID]model.Order
	mtx    sync.RWMutex
}

func New() *repository {
	return &repository{
		orders: make(map[uuid.UUID]model.Order),
	}
}

func (r *repository) SaveOrder(_ context.Context, order model.Order) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.orders[order.UUID] = model.Order{
		UUID:      order.UUID,
		Info:      order.Info,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	return nil
}

func (r *repository) GetOrder(_ context.Context, uuid uuid.UUID) (model.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	order, ok := r.orders[uuid]
	if !ok {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return order, nil
}

func (r *repository) UpdateOrder(_ context.Context, uuid uuid.UUID, info model.OrderInfo) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	order, ok := r.orders[uuid]
	if !ok {
		return errs.ErrOrderNotFound
	}

	order.Info = info
	order.UpdatedAt = lo.ToPtr(time.Now())

	r.orders[uuid] = order

	return nil
}
