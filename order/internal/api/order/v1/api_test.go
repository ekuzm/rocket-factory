package v1_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/ekuzm/rocket-factory/order/internal/api/order/v1"
	"github.com/ekuzm/rocket-factory/order/internal/api/order/v1/dto"
	serviceMock "github.com/ekuzm/rocket-factory/order/internal/api/order/v1/mock"
	errs "github.com/ekuzm/rocket-factory/order/internal/error"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	orderDto "github.com/ekuzm/rocket-factory/order/internal/service/order/dto"
	"github.com/ekuzm/rocket-factory/order/pkg/testutil"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

func TestCreateOrder(t *testing.T) {
	type args struct {
		ctx context.Context
		req *orderV1.CreateOrderRequest
	}

	tests := []struct {
		name string
		args args
		want *orderV1.CreateOrderResponse
		err  error
		mock func(service *serviceMock.OrderService, args args)
	}{
		{
			name: "ok: parses uuids, calls service, returns response",
			args: args{
				ctx: context.Background(),
				req: &orderV1.CreateOrderRequest{
					UserUUID:  testutil.TestUserUUID.String(),
					PartUuids: []string{testutil.TestPartUUID.String()},
				},
			},
			want: &orderV1.CreateOrderResponse{
				UUID:       testutil.TestOrderUUID.String(),
				TotalPrice: testutil.TestTotalPrice,
			},
			err: nil,
			mock: func(service *serviceMock.OrderService, args args) {
				summary := orderDto.Summary{
					OrderUUID:  testutil.TestOrderUUID,
					TotalPrice: testutil.TestTotalPrice,
				}

				service.On("CreateOrder", args.ctx, testutil.TestUserUUID, uuid.UUIDs{testutil.TestPartUUID}).Once().Return(summary, nil)
			},
		},
		{
			name: "invalid user uuid: returns ErrInvalidFormat, no service call",
			args: args{
				ctx: context.Background(),
				req: &orderV1.CreateOrderRequest{
					UserUUID:  uuid.Invalid.String(),
					PartUuids: []string{testutil.TestPartUUID.String()},
				},
			},
			want: nil,
			err:  errs.ErrInvalidFormat,
			mock: func(service *serviceMock.OrderService, args args) {
				service.AssertNotCalled(t, "CreateOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "invalid part uuids: returns ErrInvalidFormat, no service call",
			args: args{
				ctx: context.Background(),
				req: &orderV1.CreateOrderRequest{
					UserUUID:  testutil.TestUserUUID.String(),
					PartUuids: []string{uuid.Invalid.String()},
				},
			},
			want: nil,
			err:  errs.ErrInvalidFormat,
			mock: func(service *serviceMock.OrderService, args args) {
				service.AssertNotCalled(t, "CreateOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "service error: wrapped as create order",
			args: args{
				ctx: context.Background(),
				req: &orderV1.CreateOrderRequest{
					UserUUID:  testutil.TestUserUUID.String(),
					PartUuids: []string{testutil.TestPartUUID.String()},
				},
			},
			want: nil,
			err:  testutil.ErrService,
			mock: func(service *serviceMock.OrderService, args args) {
				service.On("CreateOrder", args.ctx, testutil.TestUserUUID, uuid.UUIDs{testutil.TestPartUUID}).Once().Return(orderDto.Summary{}, testutil.ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := serviceMock.NewOrderService(t)
			test.mock(service, test.args)

			args := api.NewAPI(service)

			got, err := args.CreateOrder(test.args.ctx, test.args.req)

			if test.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, test.want, got)
		})
	}
}

func TestGetOrder(t *testing.T) {
	type args struct {
		ctx    context.Context
		params orderV1.GetOrderParams
	}

	tests := []struct {
		name string
		args args
		want *orderV1.Order
		err  error
		mock func(service *serviceMock.OrderService, args args)
	}{
		{
			name: "ok: parses uuid, calls service, returns api order",
			args: args{
				ctx: context.Background(),
				params: orderV1.GetOrderParams{
					OrderUUID: testutil.TestOrderUUID.String(),
				},
			},
			want: dto.OrderToAPI(testutil.MakeOrder(t)),
			err:  nil,
			mock: func(service *serviceMock.OrderService, args args) {
				order := testutil.MakeOrder(t)
				service.On("GetOrder", args.ctx, testutil.TestOrderUUID).Once().Return(order, nil)
			},
		},
		{
			name: "invalid order uuid: returns ErrInvalidFormat, no service call",
			args: args{
				ctx: context.Background(),
				params: orderV1.GetOrderParams{
					OrderUUID: uuid.Invalid.String(),
				},
			},
			err: errs.ErrInvalidFormat,
			mock: func(service *serviceMock.OrderService, args args) {
				service.AssertNotCalled(t, "GetOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "service error: wrapped as get order",
			args: args{
				ctx: context.Background(),
				params: orderV1.GetOrderParams{
					OrderUUID: testutil.TestOrderUUID.String(),
				},
			},
			err: testutil.ErrService,
			mock: func(service *serviceMock.OrderService, args args) {
				service.On("GetOrder", args.ctx, testutil.TestOrderUUID).Once().Return(model.Order{}, testutil.ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := serviceMock.NewOrderService(t)
			test.mock(service, test.args)

			args := api.NewAPI(service)
			got, err := args.GetOrder(test.args.ctx, test.args.params)

			if test.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, test.want, got)
		})
	}
}

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
		mock func(service *serviceMock.OrderService, args args)
	}{
		{
			name: "ok: parses uuid, calls service, returns 204 no content",
			args: args{
				ctx: context.Background(),
				params: orderV1.CancelOrderParams{
					OrderUUID: testutil.TestOrderUUID.String(),
				},
			},
			want: &orderV1.NoContent{
				Code:    http.StatusNoContent,
				Message: testutil.TestNoContentMessage,
			},
			err: nil,
			mock: func(service *serviceMock.OrderService, args args) {
				service.On("CancelOrder", args.ctx, testutil.TestOrderUUID).Once().Return(nil)
			},
		},
		{
			name: "invalid order uuid: returns ErrInvalidFormat, no service call",
			args: args{
				ctx: context.Background(),
				params: orderV1.CancelOrderParams{
					OrderUUID: uuid.Invalid.String(),
				},
			},
			want: nil,
			err:  errs.ErrInvalidFormat,
			mock: func(service *serviceMock.OrderService, args args) {},
		},
		{
			name: "service error: wrapped as cancel order",
			args: args{
				ctx: context.Background(),
				params: orderV1.CancelOrderParams{
					OrderUUID: testutil.TestOrderUUID.String(),
				},
			},
			want: nil,
			err:  testutil.ErrService,
			mock: func(service *serviceMock.OrderService, args args) {
				service.On("CancelOrder", args.ctx, testutil.TestOrderUUID).Once().Return(testutil.ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := serviceMock.NewOrderService(t)
			test.mock(service, test.args)

			args := api.NewAPI(service)

			got, err := args.CancelOrder(test.args.ctx, test.args.params)

			if test.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, test.want, got)
		})
	}
}

func TestPayOrder(t *testing.T) {
	type args struct {
		ctx    context.Context
		req    *orderV1.PayOrderRequest
		params orderV1.PayOrderParams
	}

	tests := []struct {
		name string
		args args
		want *orderV1.PayOrderResponse
		err  error
		mock func(service *serviceMock.OrderService, args args)
	}{
		{
			name: "ok: parses uuid, maps payment method, calls service, returns transaction uuid",
			args: args{
				ctx: context.Background(),
				req: &orderV1.PayOrderRequest{
					PaymentMethod: orderV1.PaymentMethodCARD,
				},
				params: orderV1.PayOrderParams{
					OrderUUID: testutil.TestOrderUUID.String(),
				},
			},
			want: &orderV1.PayOrderResponse{
				TransactionUUID: testutil.TestTransactionUUID.String(),
			},
			err: nil,
			mock: func(service *serviceMock.OrderService, args args) {
				paymentMethod := dto.PaymentMethodToModel[args.req.PaymentMethod]

				service.On("PayOrder", args.ctx, testutil.TestOrderUUID, paymentMethod).Once().Return(testutil.TestTransactionUUID, nil)
			},
		},
		{
			name: "invalid order uuid: returns ErrInvalidFormat, no service call",
			args: args{
				ctx: context.Background(),
				req: &orderV1.PayOrderRequest{
					PaymentMethod: orderV1.PaymentMethodCARD,
				},
				params: orderV1.PayOrderParams{
					OrderUUID: uuid.Invalid.String(),
				},
			},
			want: nil,
			err:  errs.ErrInvalidFormat,
			mock: func(service *serviceMock.OrderService, args args) {},
		},
		{
			name: "service error: wrapped as pay order",
			args: args{
				ctx: context.Background(),
				req: &orderV1.PayOrderRequest{
					PaymentMethod: orderV1.PaymentMethodCARD,
				},
				params: orderV1.PayOrderParams{
					OrderUUID: testutil.TestOrderUUID.String(),
				},
			},
			want: nil,
			err:  testutil.ErrService,
			mock: func(service *serviceMock.OrderService, args args) {
				paymentMethod := dto.PaymentMethodToModel[args.req.PaymentMethod]
				service.On("PayOrder", args.ctx, testutil.TestOrderUUID, paymentMethod).Once().Return(uuid.Nil, testutil.ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			serviceMock := serviceMock.NewOrderService(t)
			test.mock(serviceMock, test.args)

			args := api.NewAPI(serviceMock)

			got, err := args.PayOrder(test.args.ctx, test.args.req, test.args.params)

			if test.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, test.want, got)
		})
	}
}

func TestNewError(t *testing.T) {
	type args struct {
		ctx context.Context
		err error
	}

	tests := []struct {
		name string
		args args
		want *orderV1.GenericErrorStatusCode
	}{
		{
			name: "returns nil",
			args: args{
				ctx: context.Background(),
				err: nil,
			},
			want: nil,
		},
		{
			name: "returns HTTP status bad request",
			args: args{
				ctx: context.Background(),
				err: errs.ErrInvalidFormat,
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusBadRequest,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusBadRequest),
					Message: orderV1.NewOptString(errs.ErrInvalidFormat.Error()),
				},
			},
		},
		{
			name: "returns HTTP status not found",
			args: args{
				ctx: context.Background(),
				err: errs.ErrOrderNotFound,
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusNotFound,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusNotFound),
					Message: orderV1.NewOptString(errs.ErrOrderNotFound.Error()),
				},
			},
		},
		{
			name: "returns HTTP status conflict by status cancelled error",
			args: args{
				ctx: context.Background(),
				err: errs.ErrStatusCancelled,
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusConflict,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusConflict),
					Message: orderV1.NewOptString(errs.ErrStatusCancelled.Error()),
				},
			},
		},
		{
			name: "returns HTTP status conflict by status paid error",
			args: args{
				ctx: context.Background(),
				err: errs.ErrStatusPaid,
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusConflict,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusConflict),
					Message: orderV1.NewOptString(errs.ErrStatusPaid.Error()),
				},
			},
		},
		{
			name: "returns HTTP status internal server error",
			args: args{
				ctx: context.Background(),
				err: testutil.ErrService,
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusInternalServerError,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusInternalServerError),
					Message: orderV1.NewOptString(testutil.ErrService.Error()),
				},
			},
		},
		{
			name: "returns HTTP status bad request from gRPC invalid argument",
			args: args{
				ctx: context.Background(),
				err: status.Error(codes.InvalidArgument, errs.ErrInvalidFormat.Error()),
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusBadRequest,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusBadRequest),
					Message: orderV1.NewOptString(errs.ErrInvalidFormat.Error()),
				},
			},
		},
		{
			name: "returns HTTP status not found from gRPC not found",
			args: args{
				ctx: context.Background(),
				err: status.Error(codes.NotFound, errs.ErrOrderNotFound.Error()),
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusNotFound,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusNotFound),
					Message: orderV1.NewOptString(errs.ErrOrderNotFound.Error()),
				},
			},
		},
		{
			name: "returns HTTP status internal server error",
			args: args{
				ctx: context.Background(),
				err: status.Error(codes.Internal, testutil.ErrService.Error()),
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusInternalServerError,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusInternalServerError),
					Message: orderV1.NewOptString(testutil.ErrService.Error()),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := serviceMock.NewOrderService(t)

			api := api.NewAPI(service)

			resp := api.NewError(test.args.ctx, test.args.err)

			if resp != nil {
				require.NotNil(t, resp)
				require.Equal(t, test.want, resp)

				return
			}

			require.Nil(t, resp)
		})
	}
}
