package main

import (
	"context"
	"errors"
	"log/slog"
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

	slog.Debug("inventory service configured", "service", "inventory-service")

	lis, err := net.Listen("tcp", config.App().GRPC.Address())
	if err != nil {
		slog.Error("inventory service listen failed", "service", "inventory-service", "addr", config.App().GRPC.Address(), "error", err)
		os.Exit(1)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil && !errors.Is(cerr, net.ErrClosed) {
			slog.Warn("inventory listener close failed", "service", "inventory-service", "error", cerr)
		}
	}()

	slog.Debug("inventory listener started", "service", "inventory-service", "addr", config.App().GRPC.Address())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mng.Connect(ctx, options.Client().ApplyURI(config.App().Mongo.URI()))
	if err != nil {
		slog.Error("mongo client create failed", "service", "inventory-service", "component", "mongo", "db", config.App().Mongo.Name(), "error", err)
		os.Exit(1)
	}
	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()

		if cerr := client.Disconnect(disconnectCtx); cerr != nil {
			slog.Warn("mongo disconnect failed", "service", "inventory-service", "component", "mongo", "error", cerr)
		}
	}()

	slog.Debug("mongo client created", "service", "inventory-service", "component", "mongo")

	if err = client.Ping(ctx, nil); err != nil {
		slog.Error("mongo ping failed", "service", "inventory-service", "component", "mongo", "error", err)
		os.Exit(1)
	}

	slog.Debug("mongo ping ok", "service", "inventory-service", "component", "mongo")

	db := client.Database(config.App().Mongo.Name())

	repository := mongo.New(db)
	service := service.New(repository)
	api := api.New(service)

	slog.Debug("inventory deps ready", "service", "inventory-service")

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.RequestLogger(), interceptor.MappingErrors()))

	inventoryV1.RegisterInventoryServiceServer(server, api)
	reflection.Register(server)

	go func() {
		slog.Debug("grpc server starting", "service", "inventory-service", "addr", config.App().GRPC.Address())
		if serveErr := server.Serve(lis); serveErr != nil && !errors.Is(serveErr, grpc.ErrServerStopped) {
			slog.Error("grpc server serve failed", "service", "inventory-service", "error", serveErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	shutdownSignal := <-quit

	slog.Warn("grpc server stopping", "service", "inventory-service", "signal", shutdownSignal.String())

	server.GracefulStop()

	slog.Debug("grpc server stopped", "service", "inventory-service")
}
