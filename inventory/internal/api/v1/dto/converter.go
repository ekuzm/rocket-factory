package dto

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/ekuzm/rocket-factory/inventory/internal/error"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartToAPI(part model.Part) *inventoryV1.Part {
	var dimensions *inventoryV1.Dimensions
	if part.Dimensions != nil {
		dimensions = &inventoryV1.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Weight: part.Dimensions.Weight,
			Height: part.Dimensions.Height,
		}
	}

	var manufacturer *inventoryV1.Manufacturer
	if part.Manufacturer != nil {
		manufacturer = &inventoryV1.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		}
	}

	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.UUID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      inventoryV1.Category(part.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          part.Tags,
		Metadata:      metadataToAPI(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     updatedAt,
	}
}

func metadataToAPI(metadata map[string]*model.Value) map[string]*inventoryV1.Value {
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

func FilterToModel(filter *inventoryV1.PartsFilter) (model.Filter, error) {
	uuids := make(uuid.UUIDs, len(filter.Uuids))

	for i, partUUID := range filter.Uuids {
		uuid, err := uuid.Parse(partUUID)
		if err != nil {
			return model.Filter{}, fmt.Errorf("parse part UUIDs: %w", errs.ErrInvalidUUIDFormat)
		}

		uuids[i] = uuid
	}

	return model.Filter{
		UUIDs:                 uuids,
		Names:                 filter.Names,
		Categories:            categoriesToModel(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}, nil
}

func categoriesToModel(categories []inventoryV1.Category) []model.Category {
	out := make([]model.Category, len(categories))

	for i, category := range categories {
		out[i] = model.Category(category)
	}

	return out
}

func PartsToAPI(parts []model.Part) []*inventoryV1.Part {
	out := make([]*inventoryV1.Part, len(parts))

	for i, part := range parts {
		out[i] = PartToAPI(part)
	}

	return out
}
