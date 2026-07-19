package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	"github.com/ekuzm/rocket-factory/order/internal/model/supplier"
	"github.com/ekuzm/rocket-factory/order/internal/service"
	dto "github.com/ekuzm/rocket-factory/order/internal/service/dto"
	mockOrder "github.com/ekuzm/rocket-factory/order/internal/service/mock"
	"github.com/ekuzm/rocket-factory/order/pkg/testutil"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
)

func TestCreateOrder(t *testing.T) {
	type args struct {
		ctx       context.Context
		userUUID  uuid.UUID
		partUUIDs uuid.UUIDs
	}

	tests := []struct {
		name string
		args args
		want dto.Summary
		err  error
		mock func(*mockOrder.OrderRepository, *mockOrder.InventoryPort, *mockOrder.TransactionManager, args)
	}{
		{
			name: "ok: sums prices and saves order",
			args: args{
				ctx:      context.Background(),
				userUUID: testutil.TestUserUUID,
				partUUIDs: uuid.UUIDs{
					testutil.TestPartUUID,
				},
			},
			want: dto.Summary{
				TotalPrice: testutil.TestTotalPrice,
			},
			err: nil,
			mock: func(repo *mockOrder.OrderRepository, port *mockOrder.InventoryPort, manager *mockOrder.TransactionManager, args args) {
				part := testutil.MakePart(t)

				port.On("ListParts", args.ctx, supplier.Filter{UUIDs: args.partUUIDs}).Once().Return([]supplier.Part{part}, nil)
				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(func(ctx context.Context, callback func(context.Context) error) error {
					return callback(ctx)
				})
				repo.On("Save", args.ctx, mock.AnythingOfType("model.Order")).Once().Return(nil)
			},
		},
		{
			name: "inventory error: returns wrapped error and does not save",
			args: args{
				ctx:      context.Background(),
				userUUID: testutil.TestUserUUID,
				partUUIDs: uuid.UUIDs{
					testutil.TestPartUUID,
				},
			},
			want: dto.Summary{},
			err:  testutil.ErrInventoryPort,
			mock: func(repo *mockOrder.OrderRepository, port *mockOrder.InventoryPort, manager *mockOrder.TransactionManager, args args) {
				port.On("ListParts", args.ctx, supplier.Filter{UUIDs: args.partUUIDs}).Once().Return(nil, testutil.ErrInventoryPort)
				manager.AssertNotCalled(t, "Wrap", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
			},
		},
		{
			name: "repository error: returns wrapped error",
			args: args{
				ctx:      context.Background(),
				userUUID: testutil.TestUserUUID,
				partUUIDs: uuid.UUIDs{
					testutil.TestPartUUID,
				},
			},
			want: dto.Summary{},
			err:  testutil.ErrRepository,
			mock: func(repo *mockOrder.OrderRepository, port *mockOrder.InventoryPort, manager *mockOrder.TransactionManager, args args) {
				part := testutil.MakePart(t)

				port.On("ListParts", args.ctx, supplier.Filter{UUIDs: args.partUUIDs}).Once().Return([]supplier.Part{part}, nil)
				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(func(ctx context.Context, callback func(context.Context) error) error {
					return callback(ctx)
				})
				repo.On("Save", args.ctx, mock.Anything).Once().Return(testutil.ErrRepository)
			},
		},
		{
			name: "ok: empty parts list => total 0, still saves",
			args: args{
				ctx:      context.Background(),
				userUUID: testutil.TestUserUUID,
				partUUIDs: uuid.UUIDs{
					testutil.TestPartUUID,
				},
			},
			want: dto.Summary{TotalPrice: 0},
			err:  nil,
			mock: func(repo *mockOrder.OrderRepository, port *mockOrder.InventoryPort, manager *mockOrder.TransactionManager, args args) {
				port.On("ListParts", args.ctx, supplier.Filter{UUIDs: args.partUUIDs}).Once().Return([]supplier.Part{}, nil)
				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(func(ctx context.Context, callback func(context.Context) error) error {
					return callback(ctx)
				})
				repo.On("Save", args.ctx, mock.MatchedBy(func(order model.Order) bool {
					return testutil.MatchOrderInput(t, order, args.userUUID, args.partUUIDs, 0)
				})).Once().Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repository := mockOrder.NewOrderRepository(t)
			inventoryPort := mockOrder.NewInventoryPort(t)
			paymentPort := mockOrder.NewPaymentPort(t)
			manager := mockOrder.NewTransactionManager(t)

			test.mock(repository, inventoryPort, manager, test.args)

			service := service.New(repository, inventoryPort, paymentPort, manager)

			got, err := service.CreateOrder(test.args.ctx, test.args.userUUID, test.args.partUUIDs)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)

				return
			}

			require.Equal(t, test.want.TotalPrice, got.TotalPrice)
		})
	}
}

func TestGetByUUID(t *testing.T) {
	type args struct {
		ctx       context.Context
		orderUUID uuid.UUID
	}

	tests := []struct {
		name string
		args args
		want model.Order
		err  error
		mock func(repo *mockOrder.OrderRepository, args args)
	}{
		{
			name: "ok: returns mapped order",
			args: args{
				ctx:       context.Background(),
				orderUUID: testutil.TestOrderUUID,
			},
			want: testutil.MakeOrder(t),
			err:  nil,
			mock: func(repo *mockOrder.OrderRepository, args args) {
				order := testutil.MakeOrder(t)

				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(order, nil)
			},
		},
		{
			name: "repository error: returns wrapped error",
			args: args{
				ctx:       context.Background(),
				orderUUID: testutil.TestOrderUUID,
			},
			want: model.Order{},
			err:  testutil.ErrRepository,
			mock: func(repo *mockOrder.OrderRepository, args args) {
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(model.Order{}, testutil.ErrRepository)
			},
		},
		{
			name: "not found: returns domain not found error",
			args: args{
				ctx:       context.Background(),
				orderUUID: testutil.TestOrderUUID,
			},
			want: model.Order{},
			err:  errs.ErrNotFound,
			mock: func(repo *mockOrder.OrderRepository, args args) {
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(model.Order{}, errs.ErrNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repository := mockOrder.NewOrderRepository(t)
			inventoryPort := mockOrder.NewInventoryPort(t)
			paymentPort := mockOrder.NewPaymentPort(t)
			manager := mockOrder.NewTransactionManager(t)

			test.mock(repository, test.args)

			svc := service.New(repository, inventoryPort, paymentPort, manager)

			got, err := svc.GetOrder(test.args.ctx, test.args.orderUUID)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Empty(t, test.want)

				return
			}

			require.Equal(t, test.want, got)
		})
	}
}

func TestCancelOrder(t *testing.T) {
	type args struct {
		ctx       context.Context
		orderUUID uuid.UUID
	}

	tests := []struct {
		name string
		args args
		err  error
		mock func(repo *mockOrder.OrderRepository, manager *mockOrder.TransactionManager, args args)
	}{
		{
			name: "ok: pending -> cancelled, updates order",
			args: args{
				ctx:       context.Background(),
				orderUUID: testutil.TestOrderUUID,
			},
			err: nil,
			mock: func(repo *mockOrder.OrderRepository, manager *mockOrder.TransactionManager, args args) {
				order := testutil.MakeOrder(t)

				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(
					func(ctx context.Context, callback func(ctx context.Context) error) error {
						return callback(ctx)
					},
				)
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(order, nil)

				order.Info.Status = model.StatusCancelled

				repo.On("Update", args.ctx, args.orderUUID, order.Info).Once().Return(nil)
			},
		},
		{
			name: "repository error on get: wrapped error",
			args: args{
				ctx:       context.Background(),
				orderUUID: testutil.TestOrderUUID,
			},
			err: testutil.ErrRepository,
			mock: func(repo *mockOrder.OrderRepository, manager *mockOrder.TransactionManager, args args) {
				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(
					func(ctx context.Context, callback func(ctx context.Context) error) error {
						return callback(ctx)
					},
				)
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(model.Order{}, testutil.ErrRepository)
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "already cancelled: returns ErrStatusCancelled, no update",
			args: args{
				ctx:       context.Background(),
				orderUUID: testutil.TestOrderUUID,
			},
			err: errs.ErrConflict,
			mock: func(repo *mockOrder.OrderRepository, manager *mockOrder.TransactionManager, args args) {
				order := testutil.MakeOrder(t)
				order.Info.Status = model.StatusCancelled

				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(
					func(ctx context.Context, callback func(ctx context.Context) error) error {
						return callback(ctx)
					},
				)
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(order, nil)
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "already paid: returns ErrStatusPaid, no update",
			args: args{
				ctx:       context.Background(),
				orderUUID: testutil.TestOrderUUID,
			},
			err: errs.ErrConflict,
			mock: func(repo *mockOrder.OrderRepository, manager *mockOrder.TransactionManager, args args) {
				order := testutil.MakeOrder(t)
				order.Info.Status = model.StatusPaid

				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(
					func(ctx context.Context, callback func(ctx context.Context) error) error {
						return callback(ctx)
					},
				)
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(order, nil)
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "repository error on update: wrapped error",
			args: args{
				ctx:       context.Background(),
				orderUUID: testutil.TestOrderUUID,
			},
			err: testutil.ErrRepository,
			mock: func(repo *mockOrder.OrderRepository, manager *mockOrder.TransactionManager, args args) {
				order := testutil.MakeOrder(t)

				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(
					func(ctx context.Context, callback func(ctx context.Context) error) error {
						return callback(ctx)
					},
				)
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(order, nil)

				order.Info.Status = model.StatusCancelled

				repo.On("Update", args.ctx, args.orderUUID, order.Info).Once().Return(testutil.ErrRepository)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repository := mockOrder.NewOrderRepository(t)
			inventoryPort := mockOrder.NewInventoryPort(t)
			paymentPort := mockOrder.NewPaymentPort(t)
			manager := mockOrder.NewTransactionManager(t)

			test.mock(repository, manager, test.args)

			svc := service.New(repository, inventoryPort, paymentPort, manager)
			err := svc.CancelOrder(test.args.ctx, test.args.orderUUID)

			if test.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestPayOrder(t *testing.T) {
	type args struct {
		ctx           context.Context
		orderUUID     uuid.UUID
		paymentMethod model.PaymentMethod
	}

	tests := []struct {
		name string
		args args
		want uuid.UUID
		err  error
		mock func(repo *mockOrder.OrderRepository, port *mockOrder.PaymentPort, manager *mockOrder.TransactionManager, args args)
	}{
		{
			name: "ok: pays, updates order with PAID + transactionUUID",
			args: args{
				ctx:           context.Background(),
				orderUUID:     testutil.TestOrderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			want: testutil.TestTransactionUUID,
			err:  nil,
			mock: func(repo *mockOrder.OrderRepository, port *mockOrder.PaymentPort, manager *mockOrder.TransactionManager, args args) {
				order := testutil.MakeOrder(t)

				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(
					func(ctx context.Context, callback func(ctx context.Context) error) error {
						return callback(ctx)
					},
				)
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(order, nil)
				port.On("PayOrder", args.ctx, args.orderUUID, args.paymentMethod).Once().Return(testutil.TestTransactionUUID, nil)

				order.Info.Status = model.StatusPaid
				order.Info.TransactionUUID = testutil.TestTransactionUUID

				repo.On("Update", args.ctx, args.orderUUID, order.Info).Once().Return(nil)
			},
		},
		{
			name: "repository error on get: wrapped error",
			args: args{
				ctx:           context.Background(),
				orderUUID:     testutil.TestOrderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			want: uuid.Nil,
			err:  testutil.ErrRepository,
			mock: func(repo *mockOrder.OrderRepository, port *mockOrder.PaymentPort, manager *mockOrder.TransactionManager, args args) {
				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(
					func(ctx context.Context, callback func(ctx context.Context) error) error {
						return callback(ctx)
					},
				)
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(model.Order{}, testutil.ErrRepository)
				port.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			},
		},
		{
			name: "already cancelled: returns ErrStatusCancelled, no payment, no update",
			args: args{
				ctx:           context.Background(),
				orderUUID:     testutil.TestOrderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			want: uuid.Nil,
			err:  errs.ErrConflict,
			mock: func(repo *mockOrder.OrderRepository, port *mockOrder.PaymentPort, manager *mockOrder.TransactionManager, args args) {
				order := testutil.MakeOrder(t)
				order.Info.Status = model.StatusCancelled

				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(
					func(ctx context.Context, callback func(ctx context.Context) error) error {
						return callback(ctx)
					},
				)
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(order, nil)
				port.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			},
		},
		{
			name: "already paid: returns ErrStatusPaid, no payment, no update",
			args: args{
				ctx:           context.Background(),
				orderUUID:     testutil.TestOrderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			want: uuid.Nil,
			err:  errs.ErrConflict,
			mock: func(repo *mockOrder.OrderRepository, port *mockOrder.PaymentPort, manager *mockOrder.TransactionManager, args args) {
				order := testutil.MakeOrder(t)
				order.Info.Status = model.StatusPaid

				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(
					func(ctx context.Context, callback func(ctx context.Context) error) error {
						return callback(ctx)
					},
				)
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(order, nil)
				port.AssertNotCalled(t, "PayOrder", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			},
		},
		{
			name: "payment port error: wrapped error, no update",
			args: args{
				ctx:           context.Background(),
				orderUUID:     testutil.TestOrderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			want: uuid.Nil,
			err:  testutil.ErrPaymentPort,
			mock: func(repo *mockOrder.OrderRepository, port *mockOrder.PaymentPort, manager *mockOrder.TransactionManager, args args) {
				order := testutil.MakeOrder(t)

				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(
					func(ctx context.Context, callback func(ctx context.Context) error) error {
						return callback(ctx)
					},
				)
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(order, nil)
				port.On("PayOrder", args.ctx, args.orderUUID, args.paymentMethod).Once().Return(uuid.Nil, testutil.ErrPaymentPort)
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			},
		},
		{
			name: "repository error on update: wrapped error",
			args: args{
				ctx:           context.Background(),
				orderUUID:     testutil.TestOrderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			want: uuid.Nil,
			err:  testutil.ErrRepository,
			mock: func(repo *mockOrder.OrderRepository, port *mockOrder.PaymentPort, manager *mockOrder.TransactionManager, args args) {
				order := testutil.MakeOrder(t)

				manager.On("Wrap", args.ctx, mock.AnythingOfType("func(context.Context) error")).Once().Return(
					func(ctx context.Context, callback func(ctx context.Context) error) error {
						return callback(ctx)
					},
				)
				repo.On("GetByUUID", args.ctx, args.orderUUID).Once().Return(order, nil)
				port.On("PayOrder", args.ctx, args.orderUUID, args.paymentMethod).Once().Return(testutil.TestTransactionUUID, nil)

				order.Info.Status = model.StatusPaid
				order.Info.TransactionUUID = testutil.TestTransactionUUID

				repo.On("Update", args.ctx, args.orderUUID, order.Info).Once().Return(testutil.ErrRepository)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repository := mockOrder.NewOrderRepository(t)
			inventoryPort := mockOrder.NewInventoryPort(t)
			paymentPort := mockOrder.NewPaymentPort(t)
			manager := mockOrder.NewTransactionManager(t)

			test.mock(repository, paymentPort, manager, test.args)

			svc := service.New(repository, inventoryPort, paymentPort, manager)

			got, err := svc.PayOrder(test.args.ctx, test.args.orderUUID, test.args.paymentMethod)

			if test.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Equal(t, uuid.Nil, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}
