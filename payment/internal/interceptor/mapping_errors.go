package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ekuzm/rocket-factory/payment/internal/model"
)

func MappingErrors() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrInvalidFormat):
				return nil, status.Error(codes.InvalidArgument, err.Error())
			default:
				return nil, status.Error(codes.Internal, err.Error())
			}
		}

		return resp, nil
	}
}
