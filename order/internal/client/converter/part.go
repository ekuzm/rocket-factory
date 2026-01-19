package converter

import (
	"time"

	"github.com/samber/lo"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

func FilterToProto(filter model.Filter) *inventoryV1.PartsFilter {
	return &inventoryV1.PartsFilter{
		Uuids:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            CategoriesToProto(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func CategoriesToProto(categories []model.Category) []inventoryV1.Category {
	out := make([]inventoryV1.Category, len(categories))

	for ind, category := range categories {
		switch category {
		case model.CategoryEngine:
			out[ind] = inventoryV1.Category_CATEGORY_ENGINE
		case model.CategoryFuel:
			out[ind] = inventoryV1.Category_CATEGORY_FUEL
		case model.CategoryPorthole:
			out[ind] = inventoryV1.Category_CATEGORY_PORTHOLE
		case model.CategoryWing:
			out[ind] = inventoryV1.Category_CATEGORY_WING
		default:
			out[ind] = inventoryV1.Category_CATEGORY_UNSPECIFIED
		}
	}

	return out
}

func PartsToModel(parts []*inventoryV1.Part) []model.Part {
	out := make([]model.Part, len(parts))

	for ind, part := range parts {
		var createdAt *time.Time
		if part.CreatedAt != nil {
			createdAt = lo.ToPtr(part.CreatedAt.AsTime())
		}

		var updatedAt *time.Time
		if part.UpdatedAt != nil {
			updatedAt = lo.ToPtr(part.UpdatedAt.AsTime())
		}

		out[ind] = model.Part{
			UUID:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Dimensions:    DimensionsToModel(part.Dimensions),
			Manufacturer:  ManufacturerToModel(part.Manufacturer),
			Category:      CategoryToModel(part.Category),
			Tags:          part.Tags,
			Metadata:      MetadataToModel(part.Metadata),
			CreatedAt:     createdAt,
			UpdateAt:      updatedAt,
		}
	}

	return out
}

func DimensionsToModel(dimensions *inventoryV1.Dimensions) *model.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToModel(manufacturer *inventoryV1.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func CategoryToModel(category inventoryV1.Category) model.Category {
	switch category {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func MetadataToModel(metadata map[string]*inventoryV1.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}

	out := make(map[string]*model.Value, len(metadata))

	for key, value := range metadata {
		if value != nil {
			switch value.Value.(type) {
			case *inventoryV1.Value_StringValue:
				out[key] = &model.Value{StringValue: lo.ToPtr(value.GetStringValue())}
			case *inventoryV1.Value_Int64Value:
				out[key] = &model.Value{Int64Value: lo.ToPtr(value.GetInt64Value())}
			case *inventoryV1.Value_DoubleValue:
				out[key] = &model.Value{DoubleValue: lo.ToPtr(value.GetDoubleValue())}
			case *inventoryV1.Value_BoolValue:
				out[key] = &model.Value{BoolValue: lo.ToPtr(value.GetBoolValue())}
			}
		}
	}

	return out
}
