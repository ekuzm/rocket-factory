package order

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/ekuzm/rocket-factory/order/internal/dto"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	repoModel "github.com/ekuzm/rocket-factory/order/internal/repository/model"
)

const (
	testOrderUUID  = "11111111-1111-1111-1111-111111111111"
	testUserUUID   = "22222222-2222-2222-2222-222222222222"
	testPartUUID   = "44444444-4444-4444-4444-444444444444"
	testTotalPrice = 15.75
	testPartPrice  = 15.75
)

var (
	ErrRepository      = errors.New("order repository error")
	ErrInventoryClient = errors.New("inventory client error")
	ErrPaymentClient   = errors.New("payment client error")
)

func makeOrder(t *testing.T, userUUID string, status model.Status) model.Order {
	t.Helper()

	return model.Order{
		UUID:     testOrderUUID,
		UserUUID: userUUID,
		PartUUIDs: []string{
			testPartUUID,
		},
		TotalPrice:    testTotalPrice,
		PaymentMethod: model.PaymentMethodCard,
		Status:        status,
	}
}

func makePayment(t *testing.T, orderUUID, userUUID string, paymentMethod model.PaymentMethod) model.Payment {
	t.Helper()

	return model.Payment{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PaymentMethod: paymentMethod,
	}
}

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
		UpdateAt:  nil,
	}
}

func matchOrderInput(t *testing.T, order repoModel.Order, input dto.CreateOrderInput, total float64) bool {
	t.Helper()

	if err := uuid.Validate(order.UUID); err != nil {
		return false
	}

	return order.UserUUID == input.UserUUID &&
		slices.Equal(order.PartUUIDs, input.PartUUIDs) &&
		order.TotalPrice == total &&
		order.Status == repoModel.StatusPendingPayment
}
