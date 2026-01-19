package order

import (
	"context"
	"fmt"

	"github.com/ekuzm/rocket-factory/order/internal/dto"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	repoConverter "github.com/ekuzm/rocket-factory/order/internal/repository/converter"
)

func (s *service) PayOrder(ctx context.Context, input dto.PayOrderInput) (dto.PayOrderOutput, error) {
	var output dto.PayOrderOutput

	order, err := s.repository.GetOrder(ctx, input.OrderUUID)
	if err != nil {
		return output, fmt.Errorf("get order: %w", err)
	}

	if order.Status == model.StatusCancelled {
		return output, fmt.Errorf("order already cancelled: %w", model.ErrConflict)
	}
	if order.Status == model.StatusPaid {
		return output, fmt.Errorf("order already paid: %w", model.ErrConflict)
	}

	payment := model.Payment{
		OrderUUID:     input.OrderUUID,
		UserUUID:      order.UserUUID,
		PaymentMethod: input.PaymentMethod,
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, payment)
	if err != nil {
		return output, err
	}

	if err = s.repository.PayOrder(ctx, repoConverter.OrderToRepoModel(order)); err != nil {
		return output, fmt.Errorf("pay order: %w", err)
	}

	return dto.PayOrderOutput{TransactionUUID: transactionUUID}, nil
}
