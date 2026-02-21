package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	api "github.com/ekuzm/rocket-factory/inventory/internal/api/v1"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, uuid uuid.UUID) (model.Part, error)
	ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error)
}

var _ api.InventoryService = (*service)(nil)

type service struct {
	inventoryRepository InventoryRepository
}

func New(repository InventoryRepository) *service {
	return &service{
		inventoryRepository: repository,
	}
}

func (s *service) GetPart(ctx context.Context, uuid uuid.UUID) (model.Part, error) {
	part, err := s.inventoryRepository.GetPart(ctx, uuid)
	if err != nil {
		return model.Part{}, fmt.Errorf("get part: %w", err)
	}

	return part, nil
}

func (s *service) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	parts, err := s.inventoryRepository.ListParts(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list parts: %w", err)
	}

	return parts, nil
}
