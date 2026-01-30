package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	serviceMock "github.com/ekuzm/rocket-factory/order/internal/service/mock"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

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
				err: model.ErrInvalidFormat,
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusBadRequest,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusBadRequest),
					Message: orderV1.NewOptString(model.ErrInvalidFormat.Error()),
				},
			},
		},
		{
			name: "returns HTTP status not found",
			args: args{
				ctx: context.Background(),
				err: model.ErrNotFound,
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusNotFound,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusNotFound),
					Message: orderV1.NewOptString(model.ErrNotFound.Error()),
				},
			},
		},
		{
			name: "returns HTTP status conflict",
			args: args{
				ctx: context.Background(),
				err: model.ErrConflict,
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusConflict,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusConflict),
					Message: orderV1.NewOptString(model.ErrConflict.Error()),
				},
			},
		},
		{
			name: "returns HTTP status internal server error",
			args: args{
				ctx: context.Background(),
				err: ErrService,
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusInternalServerError,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusInternalServerError),
					Message: orderV1.NewOptString(ErrService.Error()),
				},
			},
		},
		{
			name: "returns HTTP status bad request from gRPC invalid argument",
			args: args{
				ctx: context.Background(),
				err: status.Error(codes.InvalidArgument, model.ErrInvalidFormat.Error()),
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusBadRequest,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusBadRequest),
					Message: orderV1.NewOptString(model.ErrInvalidFormat.Error()),
				},
			},
		},
		{
			name: "returns HTTP status not found from gRPC not found",
			args: args{
				ctx: context.Background(),
				err: status.Error(codes.NotFound, model.ErrNotFound.Error()),
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusNotFound,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusNotFound),
					Message: orderV1.NewOptString(model.ErrNotFound.Error()),
				},
			},
		},
		{
			name: "returns HTTP status internal server error",
			args: args{
				ctx: context.Background(),
				err: status.Error(codes.Internal, ErrService.Error()),
			},
			want: &orderV1.GenericErrorStatusCode{
				StatusCode: http.StatusInternalServerError,
				Response: orderV1.GenericError{
					Code:    orderV1.NewOptInt(http.StatusInternalServerError),
					Message: orderV1.NewOptString(ErrService.Error()),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := serviceMock.NewOrderService(t)

			api := NewAPI(service)

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
