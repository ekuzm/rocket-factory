package service

import (
	"context"

	"github.com/ekuzm/rocket-factory/payment/internal/dto"
)

type PaymentService interface {
	PayOrder(ctx context.Context, input dto.PayOrderInput) (dto.PayOrderOutput, error)
}
