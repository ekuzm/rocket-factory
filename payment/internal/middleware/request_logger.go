package middleware

import (
	"context"
	"log"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func RequestLogger() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		timeStart := time.Now()
		method := path.Base(info.FullMethod)

		log.Printf("Start gRPC %v method", method)

		resp, err := handler(ctx, req)

		duration := time.Since(timeStart)

		if err != nil {
			status := status.Convert(err)
			log.Printf("gRPC %v method finished with %v code and had worked for %v", method, status.Code(), duration)
			return nil, err
		}

		log.Printf("gRPC %v method finished and had worked for %v", method, duration)

		return resp, nil
	}
}
