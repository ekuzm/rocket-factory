package part

import (
	"github.com/ekuzm/rocket-factory/inventory/internal/repository"
	def "github.com/ekuzm/rocket-factory/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(repository repository.InventoryRepository) *service {
	return &service{
		inventoryRepository: repository,
	}
}
