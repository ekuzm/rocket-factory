package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ekuzm/rocket-factory/inventory/internal/api/v1/dto"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
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
		logger.WithFields(logrus.Fields{
			"error": err,
			"UUID":  uuid,
		}).Warn("Failed to parse part UUID")

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
		logger.WithFields(logrus.Fields{
			"error":  err,
			"filter": filter,
		}).Warn("Failed to convert API filter to model filter")

		return nil, fmt.Errorf("convert filter to model: %w", err)
	}

	parts, err := a.service.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &inventoryV1.ListPartsResponse{Parts: dto.PartsToAPI(parts)}, nil
}
