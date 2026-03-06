package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	api "github.com/ekuzm/rocket-factory/inventory/internal/api/v1"
	errs "github.com/ekuzm/rocket-factory/inventory/internal/error"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
)

type InventoryRepository interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (model.Part, error)
	GetAllByFilter(ctx context.Context, filter model.Filter) ([]model.Part, error)
}

var _ api.InventoryService = (*service)(nil)

type service struct {
	repository InventoryRepository
}

func New(repository InventoryRepository) *service {
	return &service{
		repository: repository,
	}
}

func (s *service) GetPart(ctx context.Context, uuid uuid.UUID) (model.Part, error) {
	part, err := s.repository.GetByUUID(ctx, uuid)
	if err != nil {
		return model.Part{}, fmt.Errorf("get part: %w", err)
	}

	return part, nil
}

func (s *service) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	parts, err := s.repository.GetAllByFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list parts: %w", err)
	}

	if len(parts) != len(filter.UUIDs) && len(filter.UUIDs) > 0 {
		return nil, errs.ErrPartsNotFound
	}

	return parts, nil
}
