package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ekuzm/rocket-factory/order/internal/dto"
	"github.com/ekuzm/rocket-factory/order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, input dto.GetOrderInput) (dto.GetOrderOutput, error) {
	var output dto.GetOrderOutput

	if _, err := uuid.Parse(input.UUID); err != nil {
		return output, fmt.Errorf("order uuid: %w", model.ErrInvalidFormat)
	}

	order, err := s.repository.GetOrder(ctx, input.UUID)
	if err != nil {
		return output, fmt.Errorf("get order: %w", err)
	}

	return dto.GetOrderOutput{Order: order}, nil
}
