package dto

import "github.com/ekuzm/rocket-factory/order/internal/model"

type GetOrderInput struct {
	UUID string
}

type GetOrderOutput struct {
	Order model.Order
}
