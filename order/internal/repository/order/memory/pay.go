package memory

import (
	"context"

	repoModel "github.com/ekuzm/rocket-factory/order/internal/repository/model"
)

func (r *repository) PayOrder(ctx context.Context, order repoModel.Order) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	order.Status = repoModel.StatusPaid
	r.orders[order.UUID] = order

	return nil
}
