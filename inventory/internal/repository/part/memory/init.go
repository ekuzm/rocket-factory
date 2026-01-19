package memory

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	repoModel "github.com/ekuzm/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) InitRepository() {
	parts := generateParts()
	r.parts = make(map[string]repoModel.Part, len(parts))

	for _, part := range parts {
		r.parts[part.UUID] = part
	}
}

func generateParts() []repoModel.Part {
	parts := make([]repoModel.Part, gofakeit.Number(10, 30))

	for ind := range parts {
		parts[ind] = repoModel.Part{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.Paragraph(1, 3, 5, " "),
			Price:         gofakeit.Price(100.0, 100000.0),
			StockQuantity: gofakeit.Int64(),
			Category:      generateCategory(gofakeit.Number(0, 4)),
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          gofakeit.ProductAudience(),
			Metadata:      generateMetadata(),
			CreatedAt:     lo.ToPtr(time.Now()),
			UpdatedAt:     lo.ToPtr(time.Now()),
		}
	}

	return parts
}

func generateCategory(category int) repoModel.Category {
	switch category {
	case 1:
		return repoModel.CategoryEngine
	case 2:
		return repoModel.CategoryFuel
	case 3:
		return repoModel.CategoryPorthole
	case 4:
		return repoModel.CategoryWing
	default:
		return repoModel.CategoryUnspecified
	}
}

func generateDimensions() *repoModel.Dimensions {
	return &repoModel.Dimensions{
		Length: gofakeit.Float64(),
		Width:  gofakeit.Float64(),
		Height: gofakeit.Float64(),
		Weight: gofakeit.Float64(),
	}
}

func generateManufacturer() *repoModel.Manufacturer {
	return &repoModel.Manufacturer{
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

func generateMetadata() map[string]*repoModel.Value {
	metadata := make(map[string]*repoModel.Value)

	for index := range gofakeit.Number(0, 10) {
		key := fmt.Sprintf("%d-%s", index, gofakeit.Word())

		switch index % 4 {
		case StringValue:
			metadata[key] = &repoModel.Value{StringValue: lo.ToPtr(gofakeit.Word())}
		case Int64Value:
			metadata[key] = &repoModel.Value{Int64Value: lo.ToPtr(gofakeit.Int64())}
		case DoubleValue:
			metadata[key] = &repoModel.Value{DoubleValue: lo.ToPtr(gofakeit.Float64())}
		case BoolValue:
			metadata[key] = &repoModel.Value{BoolValue: lo.ToPtr(gofakeit.Bool())}
		}
	}

	return metadata
}
