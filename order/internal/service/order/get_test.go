package order

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	clientMock "github.com/ekuzm/rocket-factory/order/internal/client/mock"
	"github.com/ekuzm/rocket-factory/order/internal/dto"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	repositoryMock "github.com/ekuzm/rocket-factory/order/internal/repository/mock"
)

func TestGetPart(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.GetOrderInput
	}

	tests := []struct {
		name string
		args args
		want dto.GetOrderOutput
		err  error
		mock func(*repositoryMock.OrderRepository, args)
	}{
		{
			name: "returns order",
			args: args{
				ctx: context.Background(),
				input: dto.GetOrderInput{
					UUID: testOrderUUID,
				},
			},
			want: dto.GetOrderOutput{
				Order: makeOrder(t, testUserUUID, model.StatusPendingPayment),
			},
			err: nil,
			mock: func(repo *repositoryMock.OrderRepository, args args) {
				order := makeOrder(t, testUserUUID, model.StatusPendingPayment)

				repo.On("GetOrder", args.ctx, args.input.UUID).Once().Return(order, nil)
			},
		},
		{
			name: "returns invalid format error when order uuid is empty",
			args: args{
				ctx:   context.Background(),
				input: dto.GetOrderInput{},
			},
			want: dto.GetOrderOutput{},
			err:  model.ErrInvalidFormat,
			mock: func(repo *repositoryMock.OrderRepository, args args) {
				repo.AssertNotCalled(t, "GetOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns invalid format error when order uuid is invalid",
			args: args{
				ctx: context.Background(),
				input: dto.GetOrderInput{
					UUID: uuid.Invalid.String(),
				},
			},
			want: dto.GetOrderOutput{},
			err:  model.ErrInvalidFormat,
			mock: func(repo *repositoryMock.OrderRepository, args args) {
				repo.AssertNotCalled(t, "GetOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns repository error from get order method",
			args: args{
				ctx: context.Background(),
				input: dto.GetOrderInput{
					UUID: testOrderUUID,
				},
			},
			want: dto.GetOrderOutput{},
			err:  ErrRepository,
			mock: func(repo *repositoryMock.OrderRepository, args args) {
				var order model.Order

				repo.On("GetOrder", args.ctx, args.input.UUID).Once().Return(order, ErrRepository)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repositoryMock.NewOrderRepository(t)
			inventoryClient := clientMock.NewInventoryClient(t)
			paymentClient := clientMock.NewPaymentClient(t)

			service := NewService(repo, inventoryClient, paymentClient)

			test.mock(repo, test.args)

			output, err := service.GetOrder(test.args.ctx, test.args.input)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Empty(t, output)

				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, output)
		})
	}
}
