package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	"github.com/ekuzm/rocket-factory/order/internal/service"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
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

	logger.WithFields(logrus.Fields{
		"Order UUID":   order.UUID,
		"User UUID":    order.Info.UserUUID,
		"Part Count":   len(order.Info.PartUUIDs),
		"Order Status": order.Info.Status,
	}).Debug("Saved order in memory repository")

	return nil
}

func (r *repository) GetByUUID(_ context.Context, uuid uuid.UUID) (model.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	order, ok := r.orders[uuid]
	if !ok {
		logger.WithFields(logrus.Fields{
			"Order UUID": uuid,
		}).Warn("Failed to get order by UUID, not found")

		return model.Order{}, fmt.Errorf("order: %w", errs.ErrNotFound)
	}

	return order, nil
}

func (r *repository) Update(_ context.Context, uuid uuid.UUID, info model.OrderInfo) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	order, ok := r.orders[uuid]
	if !ok {
		logger.WithFields(logrus.Fields{
			"Order UUID": uuid,
		}).Warn("Failed to update order, not found")

		return fmt.Errorf("order: %w", errs.ErrNotFound)
	}

	order.Info = info
	order.UpdatedAt = lo.ToPtr(time.Now())

	r.orders[uuid] = order

	logger.WithFields(logrus.Fields{
		"Order UUID":       uuid,
		"User UUID":        info.UserUUID,
		"Order Status":     info.Status,
		"Payment Method":   info.PaymentMethod,
		"Transaction UUID": info.TransactionUUID,
	}).Debug("Updated order in memory repository")

	return nil
}
