package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/ekuzm/rocket-factory/inventory/internal/error"
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
			case errors.Is(err, errs.ErrInvalidUUIDFormat):
				return nil, status.Error(codes.InvalidArgument, err.Error())
			case errors.Is(err, errs.ErrPartNotFound):
				return nil, status.Error(codes.NotFound, err.Error())
			default:
				return nil, status.Error(codes.Internal, err.Error())
			}
		}

		return resp, nil
	}
}
