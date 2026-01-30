package v1

import (
	"errors"
	"testing"
	"time"

	"github.com/samber/lo"

	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

const (
	testPartUUID  = "44444444-4444-4444-4444-444444444444"
	testPartPrice = 15.75
)

var ErrService = errors.New("service error")

func makePart(t *testing.T) model.Part {
	t.Helper()

	return model.Part{
		UUID:          testPartUUID,
		Name:          "test-part-name",
		Description:   "test-part-description",
		Price:         testPartPrice,
		StockQuantity: 20,
		Category:      model.CategoryEngine,
		Dimensions: &model.Dimensions{
			Length: 2.4,
			Width:  1.2,
			Height: 2.4,
			Weight: 22.8,
		},
		Manufacturer: &model.Manufacturer{
			Name:    "test-manufacturer-name",
			Country: "test-manufacturer-country",
			Website: "https://test-manufacturer-website",
		},
		Tags: []string{"test-tags"},
		Metadata: map[string]*model.Value{
			"string": {
				StringValue: lo.ToPtr("string"),
			},
		},
		CreatedAt: lo.ToPtr(time.Date(2009, time.April, 12, 12, 32, 45, 9, time.Local)),
		UpdatedAt: nil,
	}
}

func makeFilter(t *testing.T, uuid string) *inventoryV1.PartsFilter {
	t.Helper()

	return &inventoryV1.PartsFilter{
		Uuids:                 []string{uuid},
		Names:                 []string{"test-part-name"},
		Categories:            []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
		ManufacturerCountries: []string{"test-country"},
		Tags:                  []string{"test-tag"},
	}
}
