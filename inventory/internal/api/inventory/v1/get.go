package v1

import (
	"context"

	"github.com/ekuzm/rocket-factory/inventory/internal/converter"
	"github.com/ekuzm/rocket-factory/inventory/internal/dto"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	input := dto.GetPartInput{
		Uuid: req.GetUuid(),
	}

	output, err := a.service.GetPart(ctx, input)
	if err != nil {
		return nil, err // error handling in mapping errors interceptor
	}

	return &inventoryV1.GetPartResponse{Part: converter.PartToAPI(output.Part)}, nil
}
