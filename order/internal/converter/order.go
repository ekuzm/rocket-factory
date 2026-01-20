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
		PaymentMethod:   orderV1.NewOptPaymentMethod(PaymentMethodToAPI(order.PaymentMethod)),
		Status:          StatusToAPI(order.Status),
	}
}

func PaymentMethodToAPI(paymentMethod model.PaymentMethod) orderV1.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethodCard:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodSPB:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodInvestorMoney:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

func StatusToAPI(status model.Status) orderV1.OrderStatus {
	switch status {
	case model.StatusPaid:
		return orderV1.OrderStatusPAID
	case model.StatusCancelled:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

func PaymentMethodToModel(paymentMethod orderV1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCard
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSPB
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}
