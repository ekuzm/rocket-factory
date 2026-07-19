package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/ekuzm/rocket-factory/payment/internal/api/v1"
	"github.com/ekuzm/rocket-factory/payment/internal/config"
	"github.com/ekuzm/rocket-factory/payment/internal/service"
	"github.com/ekuzm/rocket-factory/platform/pkg/closer"
	"github.com/ekuzm/rocket-factory/platform/pkg/grpc/health"
	"github.com/ekuzm/rocket-factory/platform/pkg/interceptor"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	"github.com/ekuzm/rocket-factory/platform/pkg/tracer"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

type app struct {
	listener net.Listener
	server   *grpc.Server
	closer   *closer.Closer
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

	return
}

func (a *app) initDeps() error {
	inits := []func() error{
		a.initConfig,
		a.initCloser,
		a.initLogger,
		a.initTracer,
		a.initListener,
		a.initServer,
	}

	for _, fn := range inits {
		if err := fn(); err != nil {
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
	cfg := config.App().GRPC

	listener, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		return fmt.Errorf("create listener: %w", err)
	}

	a.listener = listener

	a.closer.Add(listener.Close)

	return nil
}

func (a *app) initServer() error {
	service := service.New()

	api := api.New(service)

	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(interceptor.RequestLogger(), interceptor.MappingErrors()),
	)

	paymentV1.RegisterPaymentServiceServer(server, api)
	reflection.Register(server)

	health.RegisterService(server, paymentV1.PaymentService_ServiceDesc.ServiceName)

	a.server = server

	a.closer.Add(func() error {
		server.GracefulStop()

		return nil
	})

	return nil
}

func (a *app) runServer() error {
	cfg := config.App().GRPC

	errCh := make(chan error, 1)

	go func() {
		slog.Debug(
			"Start gRPC server",
			slog.String("address", cfg.Address()),
		)
		err := a.server.Serve(a.listener)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
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

	cfg := config.App().GRPC

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout())
	defer cancel()

	if err := a.closer.Close(ctx); err != nil {
		return fmt.Errorf("close deps: %w", err)
	}

	return nil
}
