package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

type PaymentService struct {
	paymentV1.UnimplementedPaymentV1ServiceServer
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (p *PaymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if err := p.validatePayOrder(req); err != nil {
		return nil, err
	}

	transactionUUID := uuid.New().String()

	log.Printf("The payment was successful, transaction uuid: %s", transactionUUID)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}

func (p *PaymentService) validatePayOrder(req *paymentV1.PayOrderRequest) error {
	if _, err := uuid.Parse(req.OrderUuid); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := uuid.Parse(req.UserUuid); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if req.PaymentMethod == paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return status.Error(codes.InvalidArgument, "payment method is unspecified")
	}

	return nil
}

const (
	serverAddress = "127.0.0.1:50052"
)

func main() {
	service := NewPaymentService()

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Printf("Failed to create listener: %v", err)
		return
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			log.Printf("Failed to close listener: %v", cerr)
		}
	}()

	server := grpc.NewServer()

	paymentV1.RegisterPaymentV1ServiceServer(server, service)
	reflection.Register(server)

	go func() {
		log.Printf("Start gRPC server at %s", serverAddress)

		if err = server.Serve(listener); err != nil {
			log.Printf("failed to serve grpc server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down the server...")
	server.GracefulStop()

	log.Printf("Server successfully stopped")
}
