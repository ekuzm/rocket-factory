package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/order/internal/dto"
	"github.com/ekuzm/rocket-factory/order/internal/service/mock"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *orderV1.CreateOrderRequest
	}

	tests := []struct {
		name string
		args args
		want *orderV1.CreateOrderResponse
		err  error
		mock func(*mock.OrderService, args)
	}{
		{
			name: "returns create order response",
			args: args{
				ctx: context.Background(),
				req: &orderV1.CreateOrderRequest{
					UserUUID:  testUserUUID,
					PartUuids: []string{testPartUUID},
				},
			},
			want: &orderV1.CreateOrderResponse{
				UUID:       testOrderUUID,
				TotalPrice: testTotalPrice,
			},
			err: nil,
			mock: func(service *mock.OrderService, args args) {
				input := dto.CreateOrderInput{
					UserUUID:  args.req.UserUUID,
					PartUUIDs: args.req.PartUuids,
				}

				output := dto.CreateOrderOutput{
					OrderUUID:  testOrderUUID,
					TotalPrice: testTotalPrice,
				}

				service.On("CreateOrder", args.ctx, input).Once().Return(output, nil)
			},
		},
		{
			name: "returns service error from create order",
			args: args{
				ctx: context.Background(),
				req: &orderV1.CreateOrderRequest{
					UserUUID:  testUserUUID,
					PartUuids: []string{testPartUUID},
				},
			},
			want: nil,
			err:  ErrService,
			mock: func(service *mock.OrderService, args args) {
				input := dto.CreateOrderInput{
					UserUUID:  args.req.UserUUID,
					PartUUIDs: args.req.PartUuids,
				}

				var output dto.CreateOrderOutput

				service.On("CreateOrder", args.ctx, input).Once().Return(output, ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := mock.NewOrderService(t)
			api := NewAPI(service)

			test.mock(service, test.args)

			resp, err := api.CreateOrder(test.args.ctx, test.args.req)
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
