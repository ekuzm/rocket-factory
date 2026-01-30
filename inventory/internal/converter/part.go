package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartToModel(part *inventoryV1.Part) model.Part {
	return model.Part{
		UUID: part.Uuid,
	}
}

func CategoryToModel(category inventoryV1.Category) model.Category {
	switch category {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

func FilterToModel(filter *inventoryV1.PartsFilter) *model.Filter {
	if filter == nil {
		return nil
	}

	return &model.Filter{
		UUIDs:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            CategoriesToModel(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func CategoriesToModel(categories []inventoryV1.Category) []model.Category {
	out := make([]model.Category, len(categories))

	for ind, category := range categories {
		out[ind] = CategoryToModel(category)
	}

	return out
}

func PartToAPI(part model.Part) *inventoryV1.Part {
	var createdAt *timestamppb.Timestamp
	if part.CreatedAt != nil {
		createdAt = timestamppb.New(*part.CreatedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToAPI(part.Category),
		Dimensions:    DimensionsToAPI(part.Dimensions),
		Manufacturer:  ManufacturerToAPI(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToAPI(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func CategoryToAPI(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

func DimensionsToAPI(dimensions *model.Dimensions) *inventoryV1.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &inventoryV1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToAPI(manufacturer *model.Manufacturer) *inventoryV1.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToAPI(metadata map[string]*model.Value) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}

	out := make(map[string]*inventoryV1.Value, len(metadata))

	for key, value := range metadata {
		switch {
		case value.StringValue != nil:
			out[key] = &inventoryV1.Value{Value: &inventoryV1.Value_StringValue{StringValue: *value.StringValue}}
		case value.Int64Value != nil:
			out[key] = &inventoryV1.Value{Value: &inventoryV1.Value_Int64Value{Int64Value: *value.Int64Value}}
		case value.DoubleValue != nil:
			out[key] = &inventoryV1.Value{Value: &inventoryV1.Value_DoubleValue{DoubleValue: *value.DoubleValue}}
		case value.BoolValue != nil:
			out[key] = &inventoryV1.Value{Value: &inventoryV1.Value_BoolValue{BoolValue: *value.BoolValue}}
		}
	}

	return out
}

func PartsToAPI(parts []model.Part) []*inventoryV1.Part {
	if parts == nil {
		return nil
	}

	out := make([]*inventoryV1.Part, len(parts))

	for ind, part := range parts {
		out[ind] = PartToAPI(part)
	}

	return out
}
