package converter

import (
	"github.com/ekuzm/rocket-factory/order/internal/model"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

func OrderToAPI(order model.Order) *orderV1.Order {
	return &orderV1.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: orderV1.NewOptString(order.TransactionUUID),
		PaymentMethod:   orderV1.NewOptPaymentMethod(orderV1.PaymentMethod(order.PaymentMethod)),
		Status:          orderV1.OrderStatus(order.Status),
	}
}
