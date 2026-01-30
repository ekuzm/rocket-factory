package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ekuzm/rocket-factory/inventory/internal/dto"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, input dto.GetPartInput) (dto.GetPartOutput, error) {
	var output dto.GetPartOutput

	if err := uuid.Validate(input.UUID); err != nil {
		return output, fmt.Errorf("part uuid: %w", model.ErrInvalidFormat)
	}

	part, err := s.inventoryRepository.GetPart(ctx, input.UUID)
	if err != nil {
		return output, fmt.Errorf("get part: %w", err)
	}

	return dto.GetPartOutput{Part: part}, nil
}
