package memory

import (
	"context"

	repoModel "github.com/ekuzm/rocket-factory/order/internal/repository/model"
)

func (r *repository) SaveOrder(ctx context.Context, order repoModel.Order) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.orders[order.UUID] = repoModel.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          repoModel.StatusPendingPayment,
	}

	return nil
}
