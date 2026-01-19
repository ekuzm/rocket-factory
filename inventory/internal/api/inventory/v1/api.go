package v1

import (
	"github.com/ekuzm/rocket-factory/inventory/internal/service"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	service service.InventoryService
}

func NewAPI(service service.InventoryService) *api {
	return &api{
		service: service,
	}
}
