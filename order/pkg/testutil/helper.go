package testutil

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	"github.com/ekuzm/rocket-factory/order/internal/model/supplier"
)

var (
	TestOrderUUID       = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	TestUserUUID        = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	TestTransactionUUID = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	TestPartUUID        = uuid.MustParse("44444444-4444-4444-4444-444444444444")
	TestCreatedAt       = time.Date(2009, time.April, 12, 12, 32, 45, 9, time.Local)
)

const (
	TestTotalPrice       = 15.75
	TestPartPrice        = 15.75
	TestNoContentMessage = "order successfully cancelled"
)

var (
	ErrService       = errors.New("order service error")
	ErrRepository    = errors.New("order repository error")
	ErrInventoryPort = errors.New("inventory port error")
	ErrPaymentPort   = errors.New("payment port error")
)

func MakeOrder(t *testing.T) model.Order {
	t.Helper()

	return model.Order{
		UUID: TestOrderUUID,
		Info: model.OrderInfo{
			UserUUID: TestUserUUID,
			PartUUIDs: uuid.UUIDs{
				TestPartUUID,
			},
			TotalPrice:    TestTotalPrice,
			PaymentMethod: model.PaymentMethodCard,
			Status:        model.StatusPendingPayment,
		},
		CreatedAt: TestCreatedAt,
		UpdatedAt: nil,
	}
}

func MakePart(t *testing.T) supplier.Part {
	t.Helper()

	return supplier.Part{
		UUID:          TestPartUUID,
		Name:          "test-part-name",
		Description:   "test-part-description",
		Price:         TestPartPrice,
		StockQuantity: 20,
		Category:      supplier.CategoryEngine,
		Dimensions: &supplier.Dimensions{
			Length: 2.4,
			Width:  1.2,
			Height: 2.4,
			Weight: 22.8,
		},
		Manufacturer: &supplier.Manufacturer{
			Name:    "test-manufacturer-name",
			Country: "test-manufacturer-country",
			Website: "https://test-manufacturer-website",
		},
		Tags: []string{"test-tags"},
		Metadata: map[string]*supplier.Value{
			"string": {
				StringValue: lo.ToPtr("string"),
			},
		},
		CreatedAt: TestCreatedAt,
		UpdateAt:  nil,
	}
}

func MatchOrderInput(t *testing.T, order model.Order, userUUID uuid.UUID, partUUIDs uuid.UUIDs, totalPrice float64) bool {
	t.Helper()

	if err := uuid.Validate(order.UUID.String()); err != nil {
		return false
	}

	return order.Info.UserUUID == userUUID &&
		slices.Equal(order.Info.PartUUIDs, partUUIDs) &&
		order.Info.TotalPrice == totalPrice &&
		order.Info.Status == model.StatusPendingPayment
}
