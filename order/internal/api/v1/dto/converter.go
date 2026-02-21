package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

func OrderToAPI(order model.Order) *orderV1.Order {
	var transactionUUID string
	if order.Info.TransactionUUID != uuid.Nil {
		transactionUUID = order.Info.TransactionUUID.String()
	}

	var updatedAt time.Time
	if order.UpdatedAt != nil {
		updatedAt = *order.UpdatedAt
	}

	return &orderV1.Order{
		UUID:            order.UUID.String(),
		UserUUID:        order.Info.UserUUID.String(),
		PartUuids:       order.Info.PartUUIDs.Strings(),
		TotalPrice:      float64(order.Info.TotalPrice),
		TransactionUUID: orderV1.NewOptString(transactionUUID),
		PaymentMethod:   orderV1.NewOptPaymentMethod(PaymentMethodToAPI[order.Info.PaymentMethod]),
		Status:          statusToAPI[order.Info.Status],
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       orderV1.OptDateTime{Value: updatedAt},
	}
}

var PaymentMethodToAPI = map[model.PaymentMethod]orderV1.PaymentMethod{
	model.PaymentMethodUnknown:       orderV1.PaymentMethodUNKNOWN,
	model.PaymentMethodSPB:           orderV1.PaymentMethodSBP,
	model.PaymentMethodCard:          orderV1.PaymentMethodCARD,
	model.PaymentMethodCreditCard:    orderV1.PaymentMethodCREDITCARD,
	model.PaymentMethodInvestorMoney: orderV1.PaymentMethodINVESTORMONEY,
}

var statusToAPI = map[model.Status]orderV1.OrderStatus{
	model.StatusPaid:           orderV1.OrderStatusPAID,
	model.StatusCancelled:      orderV1.OrderStatusCANCELLED,
	model.StatusPendingPayment: orderV1.OrderStatusPENDINGPAYMENT,
}

var PaymentMethodToModel = map[orderV1.PaymentMethod]model.PaymentMethod{
	orderV1.PaymentMethodUNKNOWN:       model.PaymentMethodUnknown,
	orderV1.PaymentMethodSBP:           model.PaymentMethodSPB,
	orderV1.PaymentMethodCARD:          model.PaymentMethodCard,
	orderV1.PaymentMethodCREDITCARD:    model.PaymentMethodCreditCard,
	orderV1.PaymentMethodINVESTORMONEY: model.PaymentMethodInvestorMoney,
}
