package converter

import (
	"github.com/ekuzm/rocket-factory/order/internal/model"
	repoModel "github.com/ekuzm/rocket-factory/order/internal/repository/model"
)

func OrderToModel(order repoModel.Order) model.Order {
	return model.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   model.PaymentMethod(order.PaymentMethod),
		Status:          model.Status(order.Status),
	}
}

func OrderToRepoModel(order model.Order) repoModel.Order {
	return repoModel.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   repoModel.PaymentMethod(order.PaymentMethod),
		Status:          repoModel.Status(order.Status),
	}
}
