package memory

import (
	"context"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	repoModel "github.com/ekuzm/rocket-factory/order/internal/repository/model"
)

func (r *repository) CancelOrder(ctx context.Context, order repoModel.Order) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	order.Status = model.StatusCancelled
	r.orders[order.UUID] = order

	return nil
}
