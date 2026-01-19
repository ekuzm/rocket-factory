package service

import (
	"context"

	"github.com/ekuzm/rocket-factory/inventory/internal/dto"
)

type InventoryService interface {
	GetPart(ctx context.Context, input dto.GetPartInput) (dto.GetPartOutput, error)
	ListParts(ctx context.Context, input dto.ListPartsInput) (dto.ListPartsOutput, error)
}
