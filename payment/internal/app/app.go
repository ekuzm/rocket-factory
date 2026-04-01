package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/ekuzm/rocket-factory/payment/internal/config"
	"github.com/ekuzm/rocket-factory/platform/pkg/closer"
	"github.com/ekuzm/rocket-factory/platform/pkg/grpc/health"
	"github.com/ekuzm/rocket-factory/platform/pkg/interceptor"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const shutdownTimeout = 5 * time.Second

type app struct {
	di       *di
	listener net.Listener
	server   *grpc.Server
}

func New(ctx context.Context) (*app, error) {
	app := &app{}

	if err := app.initDeps(ctx); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *app) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initDI,
		a.initListener,
		a.initServer,
	}

	for _, fn := range inits {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		slog.Debug("payment service starting", "addr", config.App().GRPC.Address())

		errCh <- a.runGRPCServer()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if err := closer.CloseAll(shutdownCtx); err != nil {
			return err
		}

		slog.Debug("grpc server stopped")

		return nil
	case err := <-errCh:
		return err
	}
}

func (a *app) runGRPCServer() error {
	if err := a.server.Serve(a.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("serve payment grpc server: %w", err)
	}

	return nil
}

func (a *app) initConfig(_ context.Context) error {
	return config.Setup()
}

func (a *app) initLogger(_ context.Context) error {
	logger.Init(config.App().Logger.Level(), config.App().Logger.AsJSON())

	return nil
}

func (a *app) initDI(_ context.Context) error {
	a.di = NewDI()

	return nil
}

func (a *app) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.App().GRPC.Address())
	if err != nil {
		return fmt.Errorf("create listener: %w", err)
	}

	a.listener = listener

	return nil
}

func (a *app) initServer(ctx context.Context) error {
	api, err := a.di.API(ctx)
	if err != nil {
		return fmt.Errorf("create payment api: %w", err)
	}

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.RequestLogger(), interceptor.MappingErrors()))

	paymentV1.RegisterPaymentServiceServer(server, api)
	reflection.Register(server)

	health.RegisterService(server)

	a.server = server

	closer.Add(func(ctx context.Context) error {
		server.GetServiceInfo()

		return nil
	})

	return nil
}
