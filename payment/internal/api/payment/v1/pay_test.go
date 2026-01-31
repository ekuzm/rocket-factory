package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/payment/internal/dto"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
	"github.com/ekuzm/rocket-factory/payment/internal/service/mock"
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
		mock func(*mock.PaymentService, args)
	}{
		{
			name: "returns pay order response",
			args: args{
				ctx: context.Background(),
				req: &paymentV1.PayOrderRequest{
					Uuid:          testOrderUUID,
					UserUuid:      testUserUUID,
					PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			want: &paymentV1.PayOrderResponse{
				TransactionUuid: testTransactionUUID,
			},
			err: nil,
			mock: func(service *mock.PaymentService, args args) {
				input := dto.PayOrderInput{
					OrderUUID:     args.req.Uuid,
					UserUUID:      args.req.UserUuid,
					PaymentMethod: model.PaymentMethod(args.req.PaymentMethod),
				}

				output := dto.PayOrderOutput{
					TransactionUUID: testTransactionUUID,
				}

				service.On("PayOrder", args.ctx, input).Once().Return(output, nil)
			},
		},
		{
			name: "returns service error from pay order method",
			args: args{
				ctx: context.Background(),
				req: &paymentV1.PayOrderRequest{
					Uuid:          testOrderUUID,
					UserUuid:      testUserUUID,
					PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
				},
			},
			want: &paymentV1.PayOrderResponse{
				TransactionUuid: testTransactionUUID,
			},
			err: ErrService,
			mock: func(service *mock.PaymentService, args args) {
				input := dto.PayOrderInput{
					OrderUUID:     args.req.Uuid,
					UserUUID:      args.req.UserUuid,
					PaymentMethod: model.PaymentMethod(args.req.PaymentMethod),
				}

				output := dto.PayOrderOutput{}

				service.On("PayOrder", args.ctx, input).Once().Return(output, ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := mock.NewPaymentService(t)

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
