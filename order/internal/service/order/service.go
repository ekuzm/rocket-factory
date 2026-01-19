package order

import (
	"github.com/ekuzm/rocket-factory/order/internal/client"
	"github.com/ekuzm/rocket-factory/order/internal/repository"
	def "github.com/ekuzm/rocket-factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	repository      repository.OrderRepository
	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient
}

func NewService(repository repository.OrderRepository, inventoryClient client.InventoryClient, paymentClient client.PaymentClient) *service {
	return &service{
		repository:      repository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
