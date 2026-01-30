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

func TestGetOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		params orderV1.GetOrderParams
	}

	tests := []struct {
		name string
		args args
		want *orderV1.Order
		err  error
		mock func(*mock.OrderService, args)
	}{
		{
			name: "returns get order response",
			args: args{
				ctx: context.Background(),
				params: orderV1.GetOrderParams{
					OrderUUID: testOrderUUID,
				},
			},
			want: converter.OrderToAPI(makeOrder(t)),
			err:  nil,
			mock: func(service *mock.OrderService, args args) {
				input := dto.GetOrderInput{
					UUID: args.params.OrderUUID,
				}

				output := dto.GetOrderOutput{
					Order: makeOrder(t),
				}

				service.On("GetOrder", args.ctx, input).Once().Return(output, nil)
			},
		},
		{
			name: "returns service error from get order",
			args: args{
				ctx: context.Background(),
				params: orderV1.GetOrderParams{
					OrderUUID: testOrderUUID,
				},
			},
			want: nil,
			err:  ErrService,
			mock: func(service *mock.OrderService, args args) {
				input := dto.GetOrderInput{
					UUID: args.params.OrderUUID,
				}

				var output dto.GetOrderOutput

				service.On("GetOrder", args.ctx, input).Once().Return(output, ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := mock.NewOrderService(t)
			api := NewAPI(service)

			test.mock(service, test.args)

			resp, err := api.GetOrder(test.args.ctx, test.args.params)
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
