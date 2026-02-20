package inventory

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/ekuzm/rocket-factory/order/internal/model/supplier"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

var categoryToGRPC = map[supplier.Category]inventoryV1.Category{
	supplier.CategoryUnknown:  inventoryV1.Category_CATEGORY_UNSPECIFIED,
	supplier.CategoryEngine:   inventoryV1.Category_CATEGORY_ENGINE,
	supplier.CategoryFuel:     inventoryV1.Category_CATEGORY_FUEL,
	supplier.CategoryPorthole: inventoryV1.Category_CATEGORY_PORTHOLE,
	supplier.CategoryWing:     inventoryV1.Category_CATEGORY_WING,
}

func categoriesToGRPC(categories []supplier.Category) []inventoryV1.Category {
	if categories == nil {
		return nil
	}

	out := make([]inventoryV1.Category, len(categories))

	for i, category := range categories {
		out[i] = categoryToGRPC[category]
	}

	return out
}

func filterToGRPC(filter supplier.Filter) *inventoryV1.PartsFilter {
	return &inventoryV1.PartsFilter{
		Uuids:                 filter.UUIDs.Strings(),
		Names:                 filter.Names,
		Categories:            categoriesToGRPC(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

var categoryToModel = map[inventoryV1.Category]supplier.Category{
	inventoryV1.Category_CATEGORY_UNSPECIFIED: supplier.CategoryUnknown,
	inventoryV1.Category_CATEGORY_ENGINE:      supplier.CategoryEngine,
	inventoryV1.Category_CATEGORY_FUEL:        supplier.CategoryFuel,
	inventoryV1.Category_CATEGORY_PORTHOLE:    supplier.CategoryPorthole,
	inventoryV1.Category_CATEGORY_WING:        supplier.CategoryWing,
}

func partsToModel(parts []*inventoryV1.Part) ([]supplier.Part, error) {
	if parts == nil {
		return nil, fmt.Errorf("parts are nil")
	}

	out := make([]supplier.Part, len(parts))

	for ind, part := range parts {
		var dimensions *supplier.Dimensions
		if part.Dimensions != nil {
			dimensions = &supplier.Dimensions{
				Length: part.Dimensions.Length,
				Width:  part.Dimensions.Length,
				Weight: part.Dimensions.Weight,
				Height: part.Dimensions.Height,
			}
		}

		var manufacturer *supplier.Manufacturer
		if part.Manufacturer != nil {
			manufacturer = &supplier.Manufacturer{
				Name:    part.Manufacturer.Name,
				Country: part.Manufacturer.Country,
				Website: part.Manufacturer.Website,
			}
		}

		var createdAt time.Time
		if part.CreatedAt != nil {
			createdAt = part.CreatedAt.AsTime()
		}

		var updatedAt *time.Time
		if part.UpdatedAt != nil {
			updatedAt = lo.ToPtr(part.UpdatedAt.AsTime())
		}

		uuid, err := uuid.Parse(part.Uuid)
		if err != nil {
			return nil, fmt.Errorf("uuid.Parse: %w", err)
		}

		out[ind] = supplier.Part{
			UUID:          uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Dimensions:    dimensions,
			Manufacturer:  manufacturer,
			Category:      categoryToModel[part.Category],
			Tags:          part.Tags,
			Metadata:      metadataToModel(part.Metadata),
			CreatedAt:     createdAt,
			UpdateAt:      updatedAt,
		}
	}

	return out, nil
}

func metadataToModel(metadata map[string]*inventoryV1.Value) map[string]*supplier.Value {
	if metadata == nil {
		return nil
	}

	out := make(map[string]*supplier.Value, len(metadata))

	for key, value := range metadata {
		if value != nil {
			switch value.Value.(type) {
			case *inventoryV1.Value_StringValue:
				out[key] = &supplier.Value{StringValue: lo.ToPtr(value.GetStringValue())}
			case *inventoryV1.Value_Int64Value:
				out[key] = &supplier.Value{Int64Value: lo.ToPtr(value.GetInt64Value())}
			case *inventoryV1.Value_DoubleValue:
				out[key] = &supplier.Value{DoubleValue: lo.ToPtr(value.GetDoubleValue())}
			case *inventoryV1.Value_BoolValue:
				out[key] = &supplier.Value{BoolValue: lo.ToPtr(value.GetBoolValue())}
			}
		}
	}

	return out
}
