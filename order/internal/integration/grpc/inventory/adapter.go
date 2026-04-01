package inventory

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ekuzm/rocket-factory/order/internal/model/supplier"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

type adapter struct {
	grpcClient inventoryV1.InventoryServiceClient
}

func New(grpcClient inventoryV1.InventoryServiceClient) *adapter {
	return &adapter{
		grpcClient: grpcClient,
	}
}

func (a *adapter) ListParts(ctx context.Context, filter supplier.Filter) ([]supplier.Part, error) {
	req := &inventoryV1.ListPartsRequest{
		Filter: filterToGRPC(filter),
	}

	resp, err := a.grpcClient.ListParts(ctx, req)
	if err != nil {
		slog.Error("inventory parts list failed", "error", err)

		return nil, fmt.Errorf("list parts: %w", err)
	}

	parts, err := partsToModel(resp.Parts)
	if err != nil {
		slog.Error("inventory response convert failed", "error", err)

		return nil, fmt.Errorf("parts to model: %w", err)
	}

	return parts, nil
}
