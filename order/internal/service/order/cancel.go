package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ekuzm/rocket-factory/order/internal/dto"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	repoConverter "github.com/ekuzm/rocket-factory/order/internal/repository/converter"
)

func (s *service) CancelOrder(ctx context.Context, input dto.CancelOrderInput) error {
	if _, err := uuid.Parse(input.UUID); err != nil {
		return fmt.Errorf("order uuid: %w", model.ErrInvalidFormat)
	}

	order, err := s.repository.GetOrder(ctx, input.UUID)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	if order.Status == model.StatusCancelled {
		return fmt.Errorf("order already cancelled: %w", model.ErrConflict)
	}
	if order.Status == model.StatusPaid {
		return fmt.Errorf("order already paid: %w", model.ErrConflict)
	}

	if err = s.repository.CancelOrder(ctx, repoConverter.OrderToRepoModel(order)); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	return nil
}
