package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	mng "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/ekuzm/rocket-factory/inventory/internal/api/v1"
	"github.com/ekuzm/rocket-factory/inventory/internal/config"
	"github.com/ekuzm/rocket-factory/inventory/internal/interceptor"
	"github.com/ekuzm/rocket-factory/inventory/internal/repository/mongo"
	"github.com/ekuzm/rocket-factory/inventory/internal/service"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

func main() {
	if err := config.Setup(); err != nil {
		log.Fatalf("Failed to setup config: %v", err)
	}

	lis, err := net.Listen("tcp", config.App().GRPC.Address())
	if err != nil {
		log.Fatalf("Failed to listen inventory service: %v", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil && !errors.Is(cerr, net.ErrClosed) {
			log.Printf("Failed to close inventory service listener: %v", cerr)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mng.Connect(ctx, options.Client().ApplyURI(config.App().Mongo.URI()))
	if err != nil {
		log.Printf("failed to create mongodb client: %v", err)
		return
	}
	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("failed to disconnect from mongodb: %v", cerr)
		}
	}()

	if err = client.Ping(ctx, nil); err != nil {
		log.Printf("failed to ping db: %v", err)
	}

	db := client.Database(config.App().Mongo.Name())

	repository := mongo.New(db)
	service := service.New(repository)
	api := api.New(service)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.RequestLogger(), interceptor.MappingErrors()))

	inventoryV1.RegisterInventoryServiceServer(server, api)
	reflection.Register(server)

	go func() {
		log.Printf("Start gRPC server at %s", config.App().GRPC.Address())
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
