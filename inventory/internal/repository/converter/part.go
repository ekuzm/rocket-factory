package converter

import (
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	repoModel "github.com/ekuzm/rocket-factory/inventory/internal/repository/model"
)

func PartToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      repoModel.Category(part.Category),
		Dimensions:    DimensionsToRepoModel(part.Dimensions),
		Manufacturer:  ManufacturerToRepoModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToRepoModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func DimensionsToRepoModel(dimensions *model.Dimensions) *repoModel.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &repoModel.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToRepoModel(manufacturer *model.Manufacturer) *repoModel.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToRepoModel(metadata map[string]*model.Value) map[string]*repoModel.Value {
	if metadata == nil {
		return nil
	}

	out := make(map[string]*repoModel.Value, len(metadata))

	for key, value := range metadata {
		out[key] = &repoModel.Value{
			StringValue: value.StringValue,
			Int64Value:  value.Int64Value,
			DoubleValue: value.DoubleValue,
			BoolValue:   value.BoolValue,
		}
	}

	return out
}

func FilterToRepoModel(filter *model.Filter) repoModel.Filter {
	return repoModel.Filter{
		UUIDs:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            CategoriesToRepoModel(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func CategoriesToRepoModel(categories []model.Category) []repoModel.Category {
	out := make([]repoModel.Category, len(categories))

	for ind, category := range categories {
		out[ind] = repoModel.Category(category)
	}

	return out
}

func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    DimensionsToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func DimensionsToModel(dimensions *repoModel.Dimensions) *model.Dimensions {
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

func ManufacturerToModel(manufacturer *repoModel.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToModel(metadata map[string]*repoModel.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}

	out := make(map[string]*model.Value, len(metadata))

	for key, value := range metadata {
		out[key] = &model.Value{
			StringValue: value.StringValue,
			Int64Value:  value.Int64Value,
			DoubleValue: value.DoubleValue,
			BoolValue:   value.BoolValue,
		}
	}

	return out
}

func PartsToModel(parts []repoModel.Part) []model.Part {
	out := make([]model.Part, len(parts))

	for ind, part := range parts {
		out[ind] = PartToModel(part)
	}

	return out
}
