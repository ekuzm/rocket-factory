package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	api "github.com/ekuzm/rocket-factory/inventory/internal/api/v1"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
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
		logger.WithFields(logrus.Fields{
			"Part UUID": uuid,
			"error":     err,
		}).Warn("Failed to get part by UUID")

		return model.Part{}, fmt.Errorf("get part: %w", err)
	}

	return part, nil
}

func (s *service) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	parts, err := s.repository.GetAllByFilter(ctx, filter)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Filter": filter,
			"error":  err,
		}).Error("Failed to list parts by filter")

		return nil, fmt.Errorf("list parts: %w", err)
	}

	if len(parts) != len(filter.UUIDs) && len(filter.UUIDs) > 0 {
		logger.WithFields(logrus.Fields{
			"Filter": filter,
			"Requested Part Count": len(filter.UUIDs),
			"Found Part Count":     len(parts),
		}).Warn("Failed to list parts by UUID filter")

		return nil, errs.ErrNotFound
	}

	return parts, nil
}
