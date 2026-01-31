package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/order/internal/converter"
	"github.com/ekuzm/rocket-factory/order/internal/dto"
	"github.com/ekuzm/rocket-factory/order/internal/service/mock"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		params orderV1.PayOrderParams
		req    *orderV1.PayOrderRequest
	}

	tests := []struct {
		name string
		args args
		want *orderV1.PayOrderResponse
		err  error
		mock func(*mock.OrderService, args)
	}{
		{
			name: "returns pay order response",
			args: args{
				ctx: context.Background(),
				req: &orderV1.PayOrderRequest{
					PaymentMethod: orderV1.PaymentMethodCARD,
				},
				params: orderV1.PayOrderParams{
					OrderUUID: testOrderUUID,
				},
			},
			want: &orderV1.PayOrderResponse{
				TransactionUUID: testTransactionUUID,
			},
			err: nil,
			mock: func(service *mock.OrderService, args args) {
				input := dto.PayOrderInput{
					OrderUUID:     args.params.OrderUUID,
					PaymentMethod: converter.PaymentMethodToModel(args.req.PaymentMethod),
				}

				output := dto.PayOrderOutput{
					TransactionUUID: testTransactionUUID,
				}

				service.On("PayOrder", args.ctx, input).Once().Return(output, nil)
			},
		},
		{
			name: "returns service error from pay order",
			args: args{
				ctx: context.Background(),
				req: &orderV1.PayOrderRequest{
					PaymentMethod: orderV1.PaymentMethodCARD,
				},
				params: orderV1.PayOrderParams{
					OrderUUID: testOrderUUID,
				},
			},
			want: nil,
			err:  ErrService,
			mock: func(service *mock.OrderService, args args) {
				input := dto.PayOrderInput{
					OrderUUID:     args.params.OrderUUID,
					PaymentMethod: converter.PaymentMethodToModel(args.req.PaymentMethod),
				}

				var output dto.PayOrderOutput

				service.On("PayOrder", args.ctx, input).Once().Return(output, ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := mock.NewOrderService(t)
			api := NewAPI(service)

			test.mock(service, test.args)

			resp, err := api.PayOrder(test.args.ctx, test.args.req, test.args.params)
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
