package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func RequestLogger() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			slog.Warn("grpc request finished", "method", info.FullMethod, "duration", duration, "error", err)

			return resp, err
		}

		slog.Debug("grpc request finished", "method", info.FullMethod, "duration", duration)

		return resp, nil
	}
}
