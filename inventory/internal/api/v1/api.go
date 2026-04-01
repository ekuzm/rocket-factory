package v1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/ekuzm/rocket-factory/inventory/internal/api/v1/dto"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid uuid.UUID) (model.Part, error)
	ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error)
}

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	service InventoryService
}

func New(service InventoryService) *api {
	return &api{
		service: service,
	}
}

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	uuid, err := uuid.Parse(req.Uuid)
	if err != nil {
		slog.Warn("part uuid parse failed", "uuid", req.Uuid, "error", err)

		return nil, fmt.Errorf("parse part uuid: %w", errs.ErrInvalid)
	}

	part, err := a.service.GetPart(ctx, uuid)
	if err != nil {
		return nil, err // error handling in mapping errors interceptor
	}

	return &inventoryV1.GetPartResponse{Part: dto.PartToAPI(part)}, nil
}

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter, err := dto.FilterToModel(req.Filter)
	if err != nil {
		slog.Warn(
			"api filter convert failed",
			"uuidCount", len(filter.UUIDs),
			"nameCount", len(filter.Names),
			"categoryCount", len(filter.Categories),
			"countryCount", len(filter.ManufacturerCountries),
			"tagCount", len(filter.Tags),
			"error", err,
		)

		return nil, fmt.Errorf("convert filter to model: %w", err)
	}

	parts, err := a.service.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &inventoryV1.ListPartsResponse{Parts: dto.PartsToAPI(parts)}, nil
}
