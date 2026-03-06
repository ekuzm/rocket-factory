package fixtures

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/ekuzm/rocket-factory/inventory/internal/model"
)

func GenerateParts() []model.Part {
	parts := make([]model.Part, gofakeit.Number(10, 30))

	for ind := range parts {
		parts[ind] = model.Part{
			UUID:          uuid.New(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.Paragraph(1, 3, 5, " "),
			Price:         gofakeit.Price(100.0, 100000.0),
			StockQuantity: gofakeit.Int64(),
			Category:      generateCategory(gofakeit.Number(0, 4)),
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          gofakeit.ProductAudience(),
			Metadata:      generateMetadata(),
			CreatedAt:     time.Now(),
			UpdatedAt:     lo.ToPtr(time.Now()),
		}
	}

	return parts
}

func generateCategory(category int) model.Category {
	switch category {
	case 1:
		return model.CategoryEngine
	case 2:
		return model.CategoryFuel
	case 3:
		return model.CategoryPorthole
	case 4:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

func generateDimensions() *model.Dimensions {
	return &model.Dimensions{
		Length: gofakeit.Float64(),
		Width:  gofakeit.Float64(),
		Height: gofakeit.Float64(),
		Weight: gofakeit.Float64(),
	}
}

func generateManufacturer() *model.Manufacturer {
	return &model.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: fmt.Sprintf("https://%s", gofakeit.DomainName()),
	}
}

const (
	StringValue = 0
	Int64Value  = 1
	DoubleValue = 2
	BoolValue   = 3
)

func generateMetadata() map[string]*model.Value {
	metadata := make(map[string]*model.Value)

	for index := range gofakeit.Number(0, 10) {
		key := fmt.Sprintf("%d-%s", index, gofakeit.Word())

		switch index % 4 {
		case StringValue:
			metadata[key] = &model.Value{StringValue: lo.ToPtr(gofakeit.Word())}
		case Int64Value:
			metadata[key] = &model.Value{Int64Value: lo.ToPtr(gofakeit.Int64())}
		case DoubleValue:
			metadata[key] = &model.Value{DoubleValue: lo.ToPtr(gofakeit.Float64())}
		case BoolValue:
			metadata[key] = &model.Value{BoolValue: lo.ToPtr(gofakeit.Bool())}
		}
	}

	return metadata
}
