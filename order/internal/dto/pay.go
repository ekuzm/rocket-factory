package dto

import "github.com/ekuzm/rocket-factory/order/internal/model"

type PayOrderInput struct {
	OrderUUID     string
	PaymentMethod model.PaymentMethod
}

type PayOrderOutput struct {
	TransactionUUID string
}
