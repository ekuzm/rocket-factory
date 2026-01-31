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
	repoConverter "github.com/ekuzm/rocket-factory/order/internal/repository/converter"
	repositoryMock "github.com/ekuzm/rocket-factory/order/internal/repository/mock"
)

func TestCancelOrder(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.CancelOrderInput
	}

	tests := []struct {
		name string
		args args
		err  error
		mock func(repo *repositoryMock.OrderRepository, args args)
	}{
		{
			name: "successfully cancel order",
			args: args{
				ctx: context.Background(),
				input: dto.CancelOrderInput{
					UUID: testOrderUUID,
				},
			},
			err: nil,
			mock: func(repo *repositoryMock.OrderRepository, args args) {
				order := makeOrder(t, testUserUUID, model.StatusPendingPayment)

				repo.On("GetOrder", args.ctx, args.input.UUID).Once().Return(order, nil)
				repo.On("CancelOrder", args.ctx, repoConverter.OrderToRepoModel(order)).Once().Return(nil)
			},
		},
		{
			name: "returns invalid format error when order uuid is empty",
			args: args{
				ctx:   context.Background(),
				input: dto.CancelOrderInput{},
			},
			err: model.ErrInvalidFormat,
			mock: func(repo *repositoryMock.OrderRepository, args args) {
				repo.AssertNotCalled(t, "GetOrder", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "CancelOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns invalid format error when order uuid is invalid",
			args: args{
				ctx: context.Background(),
				input: dto.CancelOrderInput{
					UUID: uuid.Invalid.String(),
				},
			},
			err: model.ErrInvalidFormat,
			mock: func(repo *repositoryMock.OrderRepository, args args) {
				repo.AssertNotCalled(t, "GetOrder", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "CancelOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns repository error from get order method",
			args: args{
				ctx: context.Background(),
				input: dto.CancelOrderInput{
					UUID: testOrderUUID,
				},
			},
			err: ErrRepository,
			mock: func(repo *repositoryMock.OrderRepository, args args) {
				var order model.Order

				repo.On("GetOrder", args.ctx, args.input.UUID).Once().Return(order, ErrRepository)
				repo.AssertNotCalled(t, "CancelOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns conflict error because order already cancelled",
			args: args{
				ctx: context.Background(),
				input: dto.CancelOrderInput{
					UUID: testOrderUUID,
				},
			},
			err: model.ErrConflict,
			mock: func(repo *repositoryMock.OrderRepository, args args) {
				order := makeOrder(t, testUserUUID, model.StatusCancelled)

				repo.On("GetOrder", args.ctx, args.input.UUID).Once().Return(order, nil)
				repo.AssertNotCalled(t, "CancelOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns conflict error because order already paid",
			args: args{
				ctx: context.Background(),
				input: dto.CancelOrderInput{
					UUID: testOrderUUID,
				},
			},
			err: model.ErrConflict,
			mock: func(repo *repositoryMock.OrderRepository, args args) {
				order := makeOrder(t, testUserUUID, model.StatusPaid)

				repo.On("GetOrder", args.ctx, args.input.UUID).Once().Return(order, nil)
				repo.AssertNotCalled(t, "CancelOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns repository error from cancel order method",
			args: args{
				ctx: context.Background(),
				input: dto.CancelOrderInput{
					UUID: testOrderUUID,
				},
			},
			err: ErrRepository,
			mock: func(repo *repositoryMock.OrderRepository, args args) {
				order := makeOrder(t, testUserUUID, model.StatusPendingPayment)

				repo.On("GetOrder", args.ctx, args.input.UUID).Once().Return(order, nil)
				repo.On("CancelOrder", args.ctx, repoConverter.OrderToRepoModel(order)).Once().Return(ErrRepository)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repo := repositoryMock.NewOrderRepository(t)
			inventoryClient := clientMock.NewInventoryClient(t)
			paymentClient := clientMock.NewPaymentClient(t)
			service := NewService(repo, inventoryClient, paymentClient)

			test.mock(repo, test.args)

			err := service.CancelOrder(test.args.ctx, test.args.input)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				return
			}

			require.NoError(t, err)
		})
	}
}
