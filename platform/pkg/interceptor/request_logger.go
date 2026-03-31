package interceptor

import (
	"context"
	"time"

	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RequestLogger() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		logger.WithFields(logrus.Fields{
			"server": info.Server,
			"method": info.FullMethod,
		}).Debug("Running gRPC method...")

		start := time.Now()

		resp, _ := handler(ctx, req)

		logger.WithFields(logrus.Fields{
			"server": info.Server,
			"method": info.FullMethod,
		}).Debug("Finished gRPC method and had worked for ", time.Since(start))

		return resp, nil
	}
}
