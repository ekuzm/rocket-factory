package memory

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	"github.com/ekuzm/rocket-factory/order/internal/service"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
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

func (r *repository) Save(_ context.Context, order model.Order) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.orders[order.UUID] = model.Order{
		UUID:      order.UUID,
		Info:      order.Info,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	slog.Debug("order saved", "orderUUID", order.UUID, "partCount", len(order.Info.PartUUIDs), "orderStatus", order.Info.Status)

	return nil
}

func (r *repository) GetByUUID(_ context.Context, uuid uuid.UUID) (model.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	order, ok := r.orders[uuid]
	if !ok {
		slog.Warn("order load missed", "orderUUID", uuid)

		return model.Order{}, fmt.Errorf("order: %w", errs.ErrNotFound)
	}

	return order, nil
}

func (r *repository) Update(_ context.Context, uuid uuid.UUID, info model.OrderInfo) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	order, ok := r.orders[uuid]
	if !ok {
		slog.Warn("order update missed", "orderUUID", uuid)

		return fmt.Errorf("order: %w", errs.ErrNotFound)
	}

	order.Info = info
	order.UpdatedAt = lo.ToPtr(time.Now())

	r.orders[uuid] = order

	slog.Debug("order updated", "orderUUID", uuid, "orderStatus", info.Status)

	return nil
}
