package memory

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	repoConverter "github.com/ekuzm/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/ekuzm/rocket-factory/inventory/internal/repository/model"
)

type PartPredicate func(repoModel.Part) bool

func (r *repository) ListParts(ctx context.Context, filter repoModel.Filter) ([]model.Part, error) {
	parts := make([]repoModel.Part, 0, len(r.parts))

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

	if len(parts) != len(filter.UUIDs) && filter.UUIDs != nil {
		return nil, fmt.Errorf("one or more parts: %w", model.ErrNotFound)
	}

	return repoConverter.PartsToModel(parts), nil
}

func filterParts(parts []repoModel.Part, preds []PartPredicate) []repoModel.Part {
	var out []repoModel.Part

	for _, part := range parts {
		if matchesAll(part, preds) {
			out = append(out, part)
		}
	}

	return out
}

func matchesAll(part repoModel.Part, preds []PartPredicate) bool {
	for _, pred := range preds {
		if !pred(part) {
			return false
		}
	}

	return true
}

func byUUIDs(uuids []string) PartPredicate {
	return func(part repoModel.Part) bool {
		for _, uuid := range uuids {
			if strings.EqualFold(part.UUID, uuid) {
				return true
			}
		}

		return false
	}
}

func byNames(names []string) PartPredicate {
	return func(part repoModel.Part) bool {
		for _, name := range names {
			if strings.EqualFold(part.Name, name) {
				return true
			}
		}

		return false
	}
}

func byCategories(categories []repoModel.Category) PartPredicate {
	return func(part repoModel.Part) bool {
		return slices.Contains(categories, part.Category)
	}
}

func byManufacturerCountries(countries []string) PartPredicate {
	return func(part repoModel.Part) bool {
		for _, country := range countries {
			if strings.EqualFold(part.Manufacturer.Country, country) {
				return true
			}
		}

		return false
	}
}

func byTags(tags []string) PartPredicate {
	return func(part repoModel.Part) bool {
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
