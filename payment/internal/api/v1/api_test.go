package v1

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockPayment "github.com/ekuzm/rocket-factory/payment/internal/api/v1/mock"
	errs "github.com/ekuzm/rocket-factory/payment/internal/error"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
	"github.com/ekuzm/rocket-factory/payment/pkg/testutil"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

func TestPayOrder(t *testing.T) {
	type args struct {
		ctx context.Context
		req *paymentV1.PayOrderRequest
	}

	tests := []struct {
		name string
		args args
		want *paymentV1.PayOrderResponse
		err  error
		mock func(*mockPayment.PaymentService, args)
	}{
		{
			name: "returns pay order response",
			args: args{
				ctx: context.Background(),
				req: &paymentV1.PayOrderRequest{
					Uuid:          testutil.TestOrderUUID.String(),
					UserUuid:      testutil.TestUserUUID.String(),
					PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			want: &paymentV1.PayOrderResponse{
				TransactionUuid: testutil.TestTransactionUUID.String(),
			},
			err: nil,
			mock: func(service *mockPayment.PaymentService, args args) {
				service.On("PayOrder", args.ctx, testutil.TestOrderUUID, testutil.TestUserUUID, model.PaymentMethodCard).Once().Return(testutil.TestTransactionUUID, nil)
			},
		},
		{
			name: "returns error for invalid order uuid",
			args: args{
				ctx: context.Background(),
				req: &paymentV1.PayOrderRequest{
					Uuid:          uuid.Invalid.String(),
					UserUuid:      testutil.TestUserUUID.String(),
					PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			want: nil,
			err:  errs.ErrInvalidUUID,
			mock: func(service *mockPayment.PaymentService, args args) {
				service.AssertNotCalled(t, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns error for invalid user uuid",
			args: args{
				ctx: context.Background(),
				req: &paymentV1.PayOrderRequest{
					Uuid:          testutil.TestOrderUUID.String(),
					UserUuid:      uuid.Invalid.String(),
					PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			want: nil,
			err:  errs.ErrInvalidUUID,
			mock: func(service *mockPayment.PaymentService, args args) {
				service.AssertNotCalled(t, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns service error from pay order method",
			args: args{
				ctx: context.Background(),
				req: &paymentV1.PayOrderRequest{
					Uuid:          testutil.TestOrderUUID.String(),
					UserUuid:      testutil.TestUserUUID.String(),
					PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
				},
			},
			want: nil,
			err:  testutil.ErrService,
			mock: func(service *mockPayment.PaymentService, args args) {
				service.On("PayOrder", args.ctx, testutil.TestOrderUUID, testutil.TestUserUUID, model.PaymentMethodUnknown).Once().Return(uuid.Nil, testutil.ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := mockPayment.NewPaymentService(t)

			api := NewAPI(service)

			test.mock(service, test.args)

			resp, err := api.PayOrder(test.args.ctx, test.args.req)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, resp)
		})
	}
}
