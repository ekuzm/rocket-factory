package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	mng "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/ekuzm/rocket-factory/inventory/internal/api/v1"
	"github.com/ekuzm/rocket-factory/inventory/internal/config"
	"github.com/ekuzm/rocket-factory/inventory/internal/repository/mongo"
	"github.com/ekuzm/rocket-factory/inventory/internal/service"
	"github.com/ekuzm/rocket-factory/platform/pkg/interceptor"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

func main() {
	if err := config.Setup(); err != nil {
		panic("Failed to setup config: " + err.Error())
	}

	logger.Init(config.App().Logger.Level(), config.App().Logger.AsJSON())

	baseLog := logger.WithFields(logrus.Fields{
		"service":        "inventory-service",
		"grpc_address":   config.App().GRPC.Address(),
		"mongo_database": config.App().Mongo.Name(),
	})

	baseLog.Debug("Loaded inventory service configuration")

	lis, err := net.Listen("tcp", config.App().GRPC.Address())
	if err != nil {
		baseLog.WithField("error", err).Fatal("Failed to listen inventory service")
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil && !errors.Is(cerr, net.ErrClosed) {
			logger.WithFields(logrus.Fields{
				"service":        "inventory-service",
				"grpc_address":   config.App().GRPC.Address(),
				"mongo_database": config.App().Mongo.Name(),
				"error":          cerr,
			}).Warn("Failed to close inventory service listener")
		}
	}()

	baseLog.Debug("Started inventory service listener")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoLog := baseLog.WithField("component", "mongo")

	client, err := mng.Connect(ctx, options.Client().ApplyURI(config.App().Mongo.URI()))
	if err != nil {
		mongoLog.WithField("error", err).Fatal("Failed to create MongoDB client")
	}
	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()

		if cerr := client.Disconnect(disconnectCtx); cerr != nil {
			mongoLog.WithField("error", cerr).Warn("Failed to disconnect from MongoDB")
		}
	}()

	mongoLog.Debug("Created MongoDB client")

	if err = client.Ping(ctx, nil); err != nil {
		mongoLog.WithField("error", err).Fatal("Failed to ping MongoDB")
	}

	mongoLog.Debug("Successfully pinged MongoDB")

	db := client.Database(config.App().Mongo.Name())

	repository := mongo.New(db)
	service := service.New(repository)
	api := api.New(service)

	logger.Debug("Inventory service dependencies initialized")

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.RequestLogger(), interceptor.MappingErrors()))

	inventoryV1.RegisterInventoryServiceServer(server, api)
	reflection.Register(server)

	go func() {
		baseLog.Debug("Starting gRPC server")
		if serveErr := server.Serve(lis); serveErr != nil && !errors.Is(serveErr, grpc.ErrServerStopped) {
			baseLog.WithField("error", serveErr).Error("Failed to serve gRPC server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	shutdownSignal := <-quit

	baseLog.WithField("signal", shutdownSignal.String()).Warn("Shutting down the gRPC server")

	server.GracefulStop()

	baseLog.Debug("gRPC server successfully stopped")
}
