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
	repoMock "github.com/ekuzm/rocket-factory/order/internal/repository/mock"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx   context.Context
		input dto.PayOrderInput
	}

	tests := []struct {
		name string
		args args
		err  error
		mock func(*repoMock.OrderRepository, *clientMock.PaymentClient, args)
	}{
		{
			name: "returns transaction UUID",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     testOrderUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: nil,
			mock: func(repo *repoMock.OrderRepository, client *clientMock.PaymentClient, args args) {
				order := makeOrder(t, testUserUUID, model.StatusPendingPayment)

				payment := makePayment(t, order.UUID, order.UserUUID, order.PaymentMethod)

				transactionUUID := uuid.NewString()

				repo.On("GetOrder", args.ctx, args.input.OrderUUID).Once().Return(order, nil)
				client.On("PayOrder", args.ctx, payment).Once().Return(transactionUUID, nil)
				repo.On("PayOrder", args.ctx, repoConverter.OrderToRepoModel(order)).Once().Return(nil)
			},
		},
		{
			name: "returns invalid format error when order uuid and payment method are invalid",
			args: args{
				ctx:   context.Background(),
				input: dto.PayOrderInput{},
			},
			err: model.ErrInvalidFormat,
			mock: func(repo *repoMock.OrderRepository, client *clientMock.PaymentClient, args args) {
				repo.AssertNotCalled(t, "GetOrder", mock.Anything, mock.Anything)
				client.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns invalid format error when order uuid is invalid",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     uuid.Invalid.String(),
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: model.ErrInvalidFormat,
			mock: func(repo *repoMock.OrderRepository, client *clientMock.PaymentClient, args args) {
				repo.AssertNotCalled(t, "GetOrder", mock.Anything, mock.Anything)
				client.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns not found error from repository get order method",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     testOrderUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: model.ErrNotFound,
			mock: func(repo *repoMock.OrderRepository, client *clientMock.PaymentClient, args args) {
				var order model.Order

				repo.On("GetOrder", args.ctx, args.input.OrderUUID).Once().Return(order, model.ErrNotFound)
				client.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns conflict error because order already cancelled",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     testOrderUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: model.ErrConflict,
			mock: func(repo *repoMock.OrderRepository, client *clientMock.PaymentClient, args args) {
				order := makeOrder(t, testUserUUID, model.StatusCancelled)

				repo.On("GetOrder", args.ctx, args.input.OrderUUID).Once().Return(order, nil)
				client.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns conflict error because order already paid",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     testOrderUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: model.ErrConflict,
			mock: func(repo *repoMock.OrderRepository, client *clientMock.PaymentClient, args args) {
				order := makeOrder(t, testUserUUID, model.StatusPaid)

				repo.On("GetOrder", args.ctx, args.input.OrderUUID).Once().Return(order, nil)
				client.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns error from payment client",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     testOrderUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: ErrPaymentClient,
			mock: func(repo *repoMock.OrderRepository, client *clientMock.PaymentClient, args args) {
				order := makeOrder(t, "", model.StatusPendingPayment)

				payment := makePayment(t, order.UUID, order.UserUUID, order.PaymentMethod)

				repo.On("GetOrder", args.ctx, args.input.OrderUUID).Once().Return(order, nil)
				client.On("PayOrder", args.ctx, payment).Once().Return("", ErrPaymentClient)
				repo.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns repository error from pay order method",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     testOrderUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: ErrRepository,
			mock: func(repo *repoMock.OrderRepository, client *clientMock.PaymentClient, args args) {
				order := makeOrder(t, testUserUUID, model.StatusPendingPayment)

				payment := makePayment(t, order.UUID, order.UserUUID, order.PaymentMethod)

				transactionUUID := uuid.NewString()

				repo.On("GetOrder", args.ctx, args.input.OrderUUID).Once().Return(order, nil)
				client.On("PayOrder", args.ctx, payment).Once().Return(transactionUUID, nil)
				repo.On("PayOrder", args.ctx, repoConverter.OrderToRepoModel(order)).Once().Return(ErrRepository)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repoMock.NewOrderRepository(t)
			inventoryClient := clientMock.NewInventoryClient(t)
			paymentClient := clientMock.NewPaymentClient(t)

			service := NewService(repo, inventoryClient, paymentClient)

			test.mock(repo, paymentClient, test.args)

			output, err := service.PayOrder(test.args.ctx, test.args.input)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Empty(t, output)

				return
			}

			require.NoError(t, err)
			require.NoError(t, uuid.Validate(output.TransactionUUID))
		})
	}
}
