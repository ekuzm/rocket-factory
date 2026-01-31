package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/order/internal/dto"
	"github.com/ekuzm/rocket-factory/order/internal/service/mock"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

func TestCancelOrder(t *testing.T) {
	type args struct {
		ctx    context.Context
		params orderV1.CancelOrderParams
	}

	tests := []struct {
		name string
		args args
		want *orderV1.NoContent
		err  error
		mock func(*mock.OrderService, args)
	}{
		{
			name: "returns cancel order response",
			args: args{
				ctx: context.Background(),
				params: orderV1.CancelOrderParams{
					OrderUUID: testOrderUUID,
				},
			},
			want: &orderV1.NoContent{
				Code:    http.StatusNoContent,
				Message: testNoContentMessage,
			},
			err: nil,
			mock: func(service *mock.OrderService, args args) {
				input := dto.CancelOrderInput{
					UUID: args.params.OrderUUID,
				}

				service.On("CancelOrder", args.ctx, input).Once().Return(nil)
			},
		},
		{
			name: "returns service error from get order method",
			args: args{
				ctx: context.Background(),
				params: orderV1.CancelOrderParams{
					OrderUUID: testOrderUUID,
				},
			},
			want: nil,
			err:  ErrService,
			mock: func(service *mock.OrderService, args args) {
				input := dto.CancelOrderInput{
					UUID: args.params.OrderUUID,
				}

				service.On("CancelOrder", args.ctx, input).Once().Return(ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := mock.NewOrderService(t)
			api := NewAPI(service)

			test.mock(service, test.args)

			resp, err := api.CancelOrder(test.args.ctx, test.args.params)
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
