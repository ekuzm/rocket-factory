package memory

import (
	"sync"

	def "github.com/ekuzm/rocket-factory/inventory/internal/repository"
	repoModel "github.com/ekuzm/rocket-factory/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	parts map[string]repoModel.Part
	mtx   sync.RWMutex
}

func NewRepository() *repository {
	repo := &repository{
		parts: make(map[string]repoModel.Part),
	}

	repo.InitRepository()

	return repo
}
