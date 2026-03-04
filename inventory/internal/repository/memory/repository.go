package memory

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"

	errs "github.com/ekuzm/rocket-factory/inventory/internal/error"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	"github.com/ekuzm/rocket-factory/inventory/internal/service"
)

var _ service.InventoryRepository = (*repository)(nil)

type repository struct {
	parts map[uuid.UUID]model.Part
	mtx   sync.RWMutex
}

func New() *repository {
	repo := &repository{
		parts: make(map[uuid.UUID]model.Part),
	}

	repo.InitRepository()

	return repo
}

func (r *repository) GetPart(_ context.Context, uuid uuid.UUID) (model.Part, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	part, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, fmt.Errorf("part with %s uuid: %w", uuid, errs.ErrPartNotFound)
	}

	return part, nil
}

type PartPredicate func(model.Part) bool

func (r *repository) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	parts := make([]model.Part, 0, len(r.parts))

	r.mtx.RLock()
	for _, part := range r.parts {
		parts = append(parts, part)
	}
	r.mtx.RUnlock()

	var preds []PartPredicate

	if len(filter.UUIDs) > 0 {
		preds = append(preds, byUUIDs(filter.UUIDs))
	}
	if len(filter.Names) > 0 {
		preds = append(preds, byNames(filter.Names))
	}
	if len(filter.Categories) > 0 {
		preds = append(preds, byCategories(filter.Categories))
	}
	if len(filter.ManufacturerCountries) > 0 {
		preds = append(preds, byManufacturerCountries(filter.ManufacturerCountries))
	}
	if len(filter.Tags) > 0 {
		preds = append(preds, byTags(filter.Tags))
	}

	parts = filterParts(parts, preds)

	if len(parts) != len(filter.UUIDs) && len(filter.UUIDs) != 0 {
		return nil, fmt.Errorf("one or more parts: %w", errs.ErrPartNotFound)
	}

	return parts, nil
}

func filterParts(parts []model.Part, preds []PartPredicate) []model.Part {
	var out []model.Part

	for _, part := range parts {
		if matchesAll(part, preds) {
			out = append(out, part)
		}
	}

	return out
}

func matchesAll(part model.Part, preds []PartPredicate) bool {
	for _, pred := range preds {
		if !pred(part) {
			return false
		}
	}

	return true
}

func byUUIDs(uuids uuid.UUIDs) PartPredicate {
	return func(part model.Part) bool {
		for _, uuid := range uuids {
			if part.UUID == uuid {
				return true
			}
		}

		return false
	}
}

func byNames(names []string) PartPredicate {
	return func(part model.Part) bool {
		for _, name := range names {
			if strings.EqualFold(part.Name, name) {
				return true
			}
		}

		return false
	}
}

func byCategories(categories []model.Category) PartPredicate {
	return func(part model.Part) bool {
		return slices.Contains(categories, part.Category)
	}
}

func byManufacturerCountries(countries []string) PartPredicate {
	return func(part model.Part) bool {
		for _, country := range countries {
			if strings.EqualFold(part.Manufacturer.Country, country) {
				return true
			}
		}

		return false
	}
}

func byTags(tags []string) PartPredicate {
	return func(part model.Part) bool {
		var counter int

		for _, partTag := range part.Tags {
			for _, tag := range tags {
				if strings.EqualFold(partTag, tag) {
					counter++
				}
			}
		}

		return counter == len(tags)
	}
}
