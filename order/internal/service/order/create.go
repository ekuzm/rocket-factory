package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ekuzm/rocket-factory/order/internal/dto"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	repoConverter "github.com/ekuzm/rocket-factory/order/internal/repository/converter"
)

func (s *service) CreateOrder(ctx context.Context, input dto.CreateOrderInput) (dto.CreateOrderOutput, error) {
	var output dto.CreateOrderOutput

	filter := model.Filter{
		UUIDs: input.PartUUIDs,
	}

	parts, err := s.inventoryClient.ListParts(ctx, filter)
	if err != nil {
		return output, err
	}

	var totalPrice float64

	for _, part := range parts {
		totalPrice += part.Price
	}

	order := model.Order{
		UUID:       uuid.New().String(),
		UserUUID:   input.UserUUID,
		PartUUIDs:  input.PartUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	if err = s.repository.SaveOrder(ctx, repoConverter.OrderToRepoModel(order)); err != nil {
		return output, fmt.Errorf("save order: %w", err)
	}

	return dto.CreateOrderOutput{OrderUUID: order.UUID, TotalPrice: order.TotalPrice}, nil
}
