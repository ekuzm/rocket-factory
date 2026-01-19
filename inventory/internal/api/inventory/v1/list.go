package v1

import (
	"context"

	"github.com/ekuzm/rocket-factory/inventory/internal/converter"
	"github.com/ekuzm/rocket-factory/inventory/internal/dto"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	input := dto.ListPartsInput{
		Filter: converter.FilterToModel(req.Filter),
	}

	output, err := a.service.ListParts(ctx, input)
	if err != nil {
		return nil, err // error handling in mapping errors interceptor
	}

	return &inventoryV1.ListPartsResponse{Parts: converter.PartsToAPI(output.Parts)}, nil
}
