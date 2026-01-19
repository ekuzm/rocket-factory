package inventory

import (
	"google.golang.org/grpc"

	def "github.com/ekuzm/rocket-factory/order/internal/client"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient inventoryV1.InventoryServiceClient
}

func NewClient(conn *grpc.ClientConn) *client {
	return &client{
		generatedClient: inventoryV1.NewInventoryServiceClient(conn),
	}
}
