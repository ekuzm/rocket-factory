package v1

import (
	"errors"
	"testing"

	"github.com/ekuzm/rocket-factory/order/internal/model"
)

const (
	testOrderUUID        = "11111111-1111-1111-1111-111111111111"
	testUserUUID         = "22222222-2222-2222-2222-222222222222"
	testPartUUID         = "44444444-4444-4444-4444-444444444444"
	testTotalPrice       = 15.75
	testTransactionUUID  = "33333333-3333-3333-3333-333333333333"
	testNoContentMessage = "order successfully cancelled"
)

var ErrService = errors.New("order service error")

func makeOrder(t *testing.T) model.Order {
	t.Helper()

	return model.Order{
		UUID:     testOrderUUID,
		UserUUID: testUserUUID,
		PartUUIDs: []string{
			testPartUUID,
		},
		TotalPrice:    testTotalPrice,
		PaymentMethod: model.PaymentMethodCard,
		Status:        model.StatusPendingPayment,
	}
}
