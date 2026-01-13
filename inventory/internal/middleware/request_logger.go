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
		method := path.Base(info.FullMethod)

		log.Printf("gRPC method %v running", method)

		timeStart := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(timeStart)

		if err != nil {
			status := status.Convert(err)
			log.Printf("Finished gRPC %v method with %v code and has worked for %v", method, status.Code(), duration)
			return nil, err
		}

		log.Printf("Finished gRPC %v method and has worked for %v", method, duration)

		return resp, nil
	}
}
