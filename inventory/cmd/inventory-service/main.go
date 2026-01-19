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

	inventoryAPI "github.com/ekuzm/rocket-factory/inventory/internal/api/inventory/v1"
	"github.com/ekuzm/rocket-factory/inventory/internal/interceptor"
	inventoryRepository "github.com/ekuzm/rocket-factory/inventory/internal/repository/part/memory"
	inventoryService "github.com/ekuzm/rocket-factory/inventory/internal/service/part"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

const (
	InventoryServiceAddress = "127.0.0.1:50051"
)

func main() {
	lis, err := net.Listen("tcp", InventoryServiceAddress)
	if err != nil {
		log.Fatalf("Failed to listen inventory service: %v", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil && !errors.Is(cerr, net.ErrClosed) {
			log.Printf("Failed to close inventory service listener: %v", cerr)
		}
	}()

	repository := inventoryRepository.NewRepository()
	service := inventoryService.NewService(repository)
	api := inventoryAPI.NewAPI(service)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.RequestLogger(), interceptor.MappingErrors()))

	inventoryV1.RegisterInventoryServiceServer(server, api)
	reflection.Register(server)

	go func() {
		log.Printf("Start gRPC server at %s", InventoryServiceAddress)
		if err := server.Serve(lis); err != nil {
			log.Printf("Failed to serve gRPC server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shuttingdown the gRPC server...")

	server.GracefulStop()

	log.Printf("gRPC server successfully stopped")
}
