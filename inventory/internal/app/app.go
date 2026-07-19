package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	mng "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/ekuzm/rocket-factory/inventory/internal/api/v1"
	"github.com/ekuzm/rocket-factory/inventory/internal/config"
	mongoRepository "github.com/ekuzm/rocket-factory/inventory/internal/repository/mongo"
	"github.com/ekuzm/rocket-factory/inventory/internal/service"
	"github.com/ekuzm/rocket-factory/platform/pkg/closer"
	"github.com/ekuzm/rocket-factory/platform/pkg/grpc/health"
	"github.com/ekuzm/rocket-factory/platform/pkg/interceptor"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	"github.com/ekuzm/rocket-factory/platform/pkg/tracer"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

type app struct {
	listener    net.Listener
	server      *grpc.Server
	mongoClient *mng.Client
	closer      *closer.Closer
}

func Run() (err error) {
	app := &app{}
	defer func() {
		err = errors.Join(err, app.closeDeps())
	}()

	if err = app.initDeps(); err != nil {
		return fmt.Errorf("init deps: %w", err)
	}

	if err = app.runServer(); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return err
}

func (a *app) initDeps() error {
	for _, init := range []func() error{
		a.initConfig,
		a.initCloser,
		a.initLogger,
		a.initTracer,
		a.initListener,
		a.initMongo,
		a.initServer,
	} {
		if err := init(); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) initConfig() error {
	return config.Setup()
}

func (a *app) initCloser() error {
	a.closer = closer.New()

	return nil
}

func (a *app) initLogger() error {
	logger.Init(config.App().Logger)

	return nil
}

func (a *app) initTracer() error {
	if err := tracer.Init(context.Background(), config.App().Tracer); err != nil {
		return fmt.Errorf("initialize tracer: %w", err)
	}

	a.closer.Add(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), config.App().GRPC.ShutdownTimeout())
		defer cancel()
		return tracer.Shutdown(ctx)
	})

	return nil
}

func (a *app) initListener() error {
	listener, err := net.Listen("tcp", config.App().GRPC.Address())
	if err != nil {
		return fmt.Errorf("create listener: %w", err)
	}

	a.listener = listener

	a.closer.Add(listener.Close)

	return nil
}

func (a *app) initMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.App().Mongo.ConnectTimeout())
	defer cancel()

	client, err := mng.Connect(ctx, options.Client().ApplyURI(config.App().Mongo.URI()).SetMonitor(otelmongo.NewMonitor()))
	if err != nil {
		return fmt.Errorf("connect mongo client: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("ping mongo client: %w", err)
	}

	a.closer.Add(func() error {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), config.App().GRPC.ShutdownTimeout())
		defer cancel()
		return client.Disconnect(shutdownCtx)
	})

	a.mongoClient = client

	return nil
}

func (a *app) initServer() error {
	repository := mongoRepository.New(context.Background(), a.mongoClient.Database(config.App().Mongo.Name()))

	inventoryService := service.New(repository)

	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(interceptor.RequestLogger(), interceptor.MappingErrors()),
	)

	inventoryV1.RegisterInventoryServiceServer(server, api.New(inventoryService))
	reflection.Register(server)

	health.RegisterService(server, inventoryV1.InventoryService_ServiceDesc.ServiceName)

	a.server = server

	a.closer.Add(func() error { server.GracefulStop(); return nil })

	return nil
}

func (a *app) runServer() error {
	errCh := make(chan error, 1)

	go func() {
		slog.Debug("Start gRPC server", slog.String("address", config.App().GRPC.Address()))
		if err := a.server.Serve(a.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errCh <- err
		}
	}()

	quit, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	select {
	case <-quit.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (a *app) closeDeps() error {
	if a.closer == nil || a.closer.IsEmpty() {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.App().GRPC.ShutdownTimeout())
	defer cancel()

	if err := a.closer.Close(ctx); err != nil {
		return fmt.Errorf("close deps: %w", err)
	}

	return nil
}
