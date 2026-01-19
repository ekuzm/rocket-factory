package inventory

import (
	"context"

	clientConverter "github.com/ekuzm/rocket-factory/order/internal/client/converter"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	req := &inventoryV1.ListPartsRequest{
		Filter: clientConverter.FilterToProto(filter),
	}

	resp, err := c.generatedClient.ListParts(ctx, req)
	if err != nil {
		return nil, err
	}

	return clientConverter.PartsToModel(resp.Parts), nil
}
