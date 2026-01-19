package part

import (
	"context"
	"fmt"

	"github.com/ekuzm/rocket-factory/inventory/internal/dto"
	"github.com/ekuzm/rocket-factory/inventory/internal/repository/converter"
)

func (s *service) ListParts(ctx context.Context, input dto.ListPartsInput) (dto.ListPartsOutput, error) {
	var output dto.ListPartsOutput

	if err := input.Filter.Validate(); err != nil {
		return output, fmt.Errorf("filter validation: %w", err)
	}

	parts, err := s.inventoryRepository.ListParts(ctx, converter.FilterToRepoModel(input.Filter))
	if err != nil {
		return output, fmt.Errorf("list parts: %w", err)
	}

	return dto.ListPartsOutput{Parts: parts}, nil
}
