package v1_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	api "github.com/ekuzm/rocket-factory/payment/internal/api/v1"
	mockPayment "github.com/ekuzm/rocket-factory/payment/internal/api/v1/mock"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
	"github.com/ekuzm/rocket-factory/payment/pkg/testutil"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
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
					PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			want: &paymentV1.PayOrderResponse{
				TransactionUuid: testutil.TestTransactionUUID.String(),
			},
			err: nil,
			mock: func(service *mockPayment.PaymentService, args args) {
				service.On("PayOrder", args.ctx, testutil.TestOrderUUID, model.PaymentMethodCard).Once().Return(testutil.TestTransactionUUID, nil)
			},
		},
		{
			name: "returns error for invalid order uuid",
			args: args{
				ctx: context.Background(),
				req: &paymentV1.PayOrderRequest{
					Uuid:          uuid.Invalid.String(),
					PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			want: nil,
			err:  errs.ErrInvalid,
			mock: func(service *mockPayment.PaymentService, args args) {
				service.AssertNotCalled(t, mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns service error from pay order method",
			args: args{
				ctx: context.Background(),
				req: &paymentV1.PayOrderRequest{
					Uuid:          testutil.TestOrderUUID.String(),
					PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
				},
			},
			want: nil,
			err:  testutil.ErrService,
			mock: func(service *mockPayment.PaymentService, args args) {
				service.On("PayOrder", args.ctx, testutil.TestOrderUUID, model.PaymentMethodUnknown).Once().Return(uuid.Nil, testutil.ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := mockPayment.NewPaymentService(t)

			api := api.New(service)

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
