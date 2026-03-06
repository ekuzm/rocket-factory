package entity

import (
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
)

var categoryToModel = map[Category]model.Category{
	CategoryUnspecified: model.CategoryUnspecified,
	CategoryEngine:      model.CategoryEngine,
	CategoryFuel:        model.CategoryFuel,
	CategoryPorthole:    model.CategoryPorthole,
	CategoryWing:        model.CategoryWing,
}

var categoryToDocument = map[model.Category]Category{
	model.CategoryUnspecified: CategoryUnspecified,
	model.CategoryEngine:      CategoryEngine,
	model.CategoryFuel:        CategoryFuel,
	model.CategoryPorthole:    CategoryPorthole,
	model.CategoryWing:        CategoryWing,
}

func metadataToModel(metadata map[string]*Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}

	out := make(map[string]*model.Value, len(metadata))

	for k, v := range metadata {
		out[k] = &model.Value{
			StringValue: v.StringValue,
			Int64Value:  v.Int64Value,
			DoubleValue: v.DoubleValue,
			BoolValue:   v.BoolValue,
		}
	}

	return out
}

func metadataToDocument(metadata map[string]*model.Value) map[string]*Value {
	if metadata == nil {
		return nil
	}

	out := make(map[string]*Value, len(metadata))

	for k, v := range metadata {
		out[k] = &Value{
			StringValue: v.StringValue,
			Int64Value:  v.Int64Value,
			DoubleValue: v.DoubleValue,
			BoolValue:   v.BoolValue,
		}
	}

	return out
}

func PartToModel(part PartDocument) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      categoryToModel[part.Category],
		Dimensions: &model.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: &model.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      part.Tags,
		Metadata:  metadataToModel(part.Metadata),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

func PartToDocument(part model.Part) PartDocument {
	return PartDocument{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      categoryToDocument[part.Category],
		Dimensions: &Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: &Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      part.Tags,
		Metadata:  metadataToDocument(part.Metadata),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

func PartsToModel(parts []PartDocument) []model.Part {
	if parts == nil {
		return nil
	}

	out := make([]model.Part, len(parts))

	for i, v := range parts {
		out[i] = PartToModel(v)
	}

	return out
}

func PartsToDocument(parts []model.Part) []PartDocument {
	if parts == nil {
		return nil
	}

	out := make([]PartDocument, len(parts))

	for i, v := range parts {
		out[i] = PartToDocument(v)
	}

	return out
}

func CategoriesToDocument(categories []model.Category) []Category {
	out := make([]Category, len(categories))

	for i, v := range categories {
		out[i] = categoryToDocument[v]
	}

	return out
}
