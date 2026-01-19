package repository

import (
	"context"

	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	repoModel "github.com/ekuzm/rocket-factory/inventory/internal/repository/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter repoModel.Filter) ([]model.Part, error)
}
