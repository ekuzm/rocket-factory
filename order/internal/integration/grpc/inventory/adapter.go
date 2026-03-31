package inventory

import (
	"context"
	"fmt"

	"github.com/ekuzm/rocket-factory/order/internal/model/supplier"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
	"github.com/sirupsen/logrus"
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
		logger.WithFields(logrus.Fields{
			"Filter": filter,
			"error":  err,
		}).Error("Failed to list parts in inventory service")

		return nil, fmt.Errorf("list parts: %w", err)
	}

	parts, err := partsToModel(resp.Parts)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Filter":     filter,
			"Part Count": len(resp.Parts),
			"error":      err,
		}).Error("Failed to convert inventory parts response")

		return nil, fmt.Errorf("parts to model: %w", err)
	}

	return parts, nil
}
