package dto

import (
	"github.com/ekuzm/rocket-factory/payment/internal/model"
)

type PayOrderInput struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod model.PaymentMethod
}

type PayOrderOutput struct {
	TransactionUUID string
}
