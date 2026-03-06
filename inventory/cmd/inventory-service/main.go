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

	"go.mongodb.org/mongo-driver/bson"
	mng "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/ekuzm/rocket-factory/inventory/internal/api/v1"
	"github.com/ekuzm/rocket-factory/inventory/internal/interceptor"
	"github.com/ekuzm/rocket-factory/inventory/internal/repository/mongo"
	"github.com/ekuzm/rocket-factory/inventory/internal/service"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

const (
	inventoryServiceAddress = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", inventoryServiceAddress)
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

	dbURI := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")

	client, err := mng.Connect(ctx, options.Client().ApplyURI(dbURI))
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

	db := client.Database(dbName)

	collection := db.Collection("parts")

	indexModel := mng.IndexModel{
		Keys: bson.D{
			{Key: "uuid", Value: 1},
			{Key: "name", Value: 1},
			{Key: "category", Value: 1},
			{Key: "manufacturer.country", Value: 1},
			{Key: "tags", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	indexName, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("failed to create index: %v", err)
	}

	log.Printf("Create index with name: %v", indexName)

	repository := mongo.New(collection)
	service := service.New(repository)
	api := api.New(service)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.RequestLogger(), interceptor.MappingErrors()))

	inventoryV1.RegisterInventoryServiceServer(server, api)
	reflection.Register(server)

	go func() {
		log.Printf("Start gRPC server at %s", inventoryServiceAddress)
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
