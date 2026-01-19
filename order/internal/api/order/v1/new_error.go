package v1

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) NewError(ctx context.Context, err error) (r *orderV1.GenericErrorStatusCode) {
	var (
		code    int
		message string
	)

	status, ok := status.FromError(err)
	if !ok {
		switch {
		case errors.Is(err, model.ErrInvalidFormat):
			code, message = http.StatusBadRequest, err.Error()
		case errors.Is(err, model.ErrNotFound):
			code, message = http.StatusNotFound, err.Error()
		case errors.Is(err, model.ErrConflict):
			code, message = http.StatusConflict, err.Error()
		default:
			code, message = http.StatusInternalServerError, err.Error()
		}
	} else {
		switch status.Code() {
		case codes.InvalidArgument:
			code, message = http.StatusBadRequest, status.Message()
		case codes.NotFound:
			code, message = http.StatusNotFound, status.Message()
		default:
			code, message = http.StatusInternalServerError, status.Message()
		}
	}

	return &orderV1.GenericErrorStatusCode{
		StatusCode: code,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(code),
			Message: orderV1.NewOptString(message),
		},
	}
}
