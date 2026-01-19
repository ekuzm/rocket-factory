package main

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentAPI "github.com/ekuzm/rocket-factory/payment/internal/api/payment/v1"
	"github.com/ekuzm/rocket-factory/payment/internal/interceptor"
	paymentService "github.com/ekuzm/rocket-factory/payment/internal/service/payment"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	PaymentServiceAddress = "127.0.0.1:50052"
)

func main() {
	lis, err := net.Listen("tcp", PaymentServiceAddress)
	if err != nil {
		log.Fatalf("Failed to listen payment service at %s: %v", PaymentServiceAddress, err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil && !errors.Is(cerr, net.ErrClosed) {
			log.Printf("failed to close listener: %v", cerr)
		}
	}()

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.RequestLogger(), interceptor.MappingErrors()))

	service := paymentService.NewService()
	api := paymentAPI.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(server, api)
	reflection.Register(server)

	go func() {
		log.Printf("Start gRPC server at %s", PaymentServiceAddress)

		if err = server.Serve(lis); err != nil {
			log.Printf("Failed to serve gRPC server at %s: %v", PaymentServiceAddress, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Printf("Shutting down the server...")

	server.GracefulStop()

	log.Printf("Server successfully stopped")
}
