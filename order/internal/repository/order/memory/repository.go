package memory

import (
	"sync"

	def "github.com/ekuzm/rocket-factory/order/internal/repository"
	repoModel "github.com/ekuzm/rocket-factory/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	orders map[string]repoModel.Order
	mtx    sync.RWMutex
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]repoModel.Order),
	}
}
