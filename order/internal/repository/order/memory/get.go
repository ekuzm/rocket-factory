package memory

import (
	"context"
	"fmt"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	repoConverter "github.com/ekuzm/rocket-factory/order/internal/repository/converter"
)

func (r *repository) GetOrder(ctx context.Context, uuid string) (model.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	order, ok := r.orders[uuid]
	if !ok {
		return model.Order{}, fmt.Errorf("order with %s uuid: %w", uuid, model.ErrNotFound)
	}

	return repoConverter.OrderToModel(order), nil
}
