package memory

import (
	"context"
	"fmt"

	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	repoConverter "github.com/ekuzm/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, uuid string) (model.Part, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	part, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, fmt.Errorf("part with %s uuid: %w", uuid, model.ErrNotFound)
	}

	return repoConverter.PartToModel(part), nil
}
