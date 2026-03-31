package main

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/ekuzm/rocket-factory/payment/internal/api/v1"
	"github.com/ekuzm/rocket-factory/payment/internal/config"
	"github.com/ekuzm/rocket-factory/payment/internal/service"
	"github.com/ekuzm/rocket-factory/platform/pkg/interceptor"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

func main() {
	if err := config.Setup(); err != nil {
		panic("Failed to setup config")
	}

	lis, err := net.Listen("tcp", config.App().GRPC.Address())
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service":      "payment-service",
			"grpc_address": config.App().GRPC.Address(),
			"error":        err,
		}).Fatal("Failed to listen payment service")
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil && !errors.Is(cerr, net.ErrClosed) {
			logger.WithFields(logrus.Fields{
				"service":      "payment-service",
				"grpc_address": config.App().GRPC.Address(),
				"error":        cerr,
			}).Warn("Failed to close payment service listener")
		}
	}()

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.RequestLogger(), interceptor.MappingErrors()))

	service := service.New()
	api := api.New(service)

	paymentV1.RegisterPaymentServiceServer(server, api)
	reflection.Register(server)

	go func() {
		logger.WithFields(logrus.Fields{
			"service":      "payment-service",
			"grpc_address": config.App().GRPC.Address(),
		}).Debug("Starting gRPC server...")

		if err := server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.WithFields(logrus.Fields{
				"service":      "payment-service",
				"grpc_address": config.App().GRPC.Address(),
				"error":        err,
			}).Error("Failed to serve gRPC server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	shutdownSignal := <-quit

	logger.WithFields(logrus.Fields{
		"service":      "payment-service",
		"grpc_address": config.App().GRPC.Address(),
		"signal":       shutdownSignal.String(),
	}).Warn("Shutting down the gRPC server")

	server.GracefulStop()

	logger.WithFields(logrus.Fields{
		"service":      "payment-service",
		"grpc_address": config.App().GRPC.Address(),
	}).Debug("gRPC server successfully stopped")
}
