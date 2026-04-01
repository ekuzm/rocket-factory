package main

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/ekuzm/rocket-factory/payment/internal/api/v1"
	"github.com/ekuzm/rocket-factory/payment/internal/config"
	"github.com/ekuzm/rocket-factory/payment/internal/service"
	"github.com/ekuzm/rocket-factory/platform/pkg/interceptor"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

func main() {
	if err := config.Setup(); err != nil {
		panic("Failed to setup config")
	}

	lis, err := net.Listen("tcp", config.App().GRPC.Address())
	if err != nil {
		slog.Error("payment service listen failed", "service", "payment-service", "addr", config.App().GRPC.Address(), "error", err)
		os.Exit(1)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil && !errors.Is(cerr, net.ErrClosed) {
			slog.Warn("payment listener close failed", "service", "payment-service", "error", cerr)
		}
	}()

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.RequestLogger(), interceptor.MappingErrors()))

	service := service.New()
	api := api.New(service)

	paymentV1.RegisterPaymentServiceServer(server, api)
	reflection.Register(server)

	go func() {
		slog.Debug("grpc server starting", "service", "payment-service", "addr", config.App().GRPC.Address())

		if err := server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("grpc server serve failed", "service", "payment-service", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	shutdownSignal := <-quit

	slog.Warn("grpc server stopping", "service", "payment-service", "signal", shutdownSignal.String())

	server.GracefulStop()

	slog.Debug("grpc server stopped", "service", "payment-service")
}
