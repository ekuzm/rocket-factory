package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/exaring/otelpgx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/ekuzm/rocket-factory/order/internal/api/v1"
	"github.com/ekuzm/rocket-factory/order/internal/config"
	"github.com/ekuzm/rocket-factory/order/internal/integration/grpc/inventory"
	"github.com/ekuzm/rocket-factory/order/internal/integration/grpc/payment"
	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres"
	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres/transaction"
	"github.com/ekuzm/rocket-factory/order/internal/service"
	"github.com/ekuzm/rocket-factory/platform/pkg/closer"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	customMiddleware "github.com/ekuzm/rocket-factory/platform/pkg/middleware"
	"github.com/ekuzm/rocket-factory/platform/pkg/tracer"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

type app struct {
	server        *http.Server
	pool          *pgxpool.Pool
	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn
	closer        *closer.Closer
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
	for _, init := range []func() error{
		a.initConfig,
		a.initCloser,
		a.initLogger,
		a.initTracer,
		a.initPool,
		a.initInventoryClient,
		a.initPaymentClient,
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
		ctx, cancel := context.WithTimeout(context.Background(), config.App().HTTP.ShutdownTimeout())
		defer cancel()
		return tracer.Shutdown(ctx)
	})

	return nil
}

func (a *app) initPool() error {
	ctx := context.Background()

	poolConfig, err := pgxpool.ParseConfig(config.App().Postgres.URI())
	if err != nil {
		return fmt.Errorf("parse database pool config: %w", err)
	}

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("create database pool: %w", err)
	}

	a.closer.Add(func() error {
		pool.Close()

		return nil
	})

	a.pool = pool

	return nil
}

func (a *app) initInventoryClient() error {
	conn, err := grpc.NewClient(
		config.App().GRPC.InventoryAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return fmt.Errorf("create inventory client connection: %w", err)
	}

	a.closer.Add(conn.Close)

	a.inventoryConn = conn

	return nil
}

func (a *app) initPaymentClient() error {
	conn, err := grpc.NewClient(
		config.App().GRPC.PaymentAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return fmt.Errorf("create payment client connection: %w", err)
	}

	a.closer.Add(conn.Close)

	a.paymentConn = conn

	return nil
}

func (a *app) initServer() error {
	orderService := service.New(
		postgres.New(a.pool),
		inventory.New(inventoryV1.NewInventoryServiceClient(a.inventoryConn)),
		payment.New(paymentV1.NewPaymentServiceClient(a.paymentConn)),
		transaction.New(a.pool),
	)

	orderServer, err := orderV1.NewServer(api.New(orderService))
	if err != nil {
		return fmt.Errorf("create order server: %w", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recoverer, middleware.Timeout(config.App().HTTP.RequestTimeout()), customMiddleware.RequestLogger)
	router.Mount("/", orderServer)

	a.server = &http.Server{
		Addr:              config.App().HTTP.OrderAddress(),
		Handler:           otelhttp.NewHandler(router, "order.http"),
		ReadHeaderTimeout: config.App().HTTP.ReadHeaderTimeout(),
	}

	a.closer.Add(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), config.App().HTTP.ShutdownTimeout())
		defer cancel()
		return a.server.Shutdown(ctx)
	})

	return nil
}

func (a *app) runServer() error {
	errCh := make(chan error, 1)

	go func() {
		slog.Debug("Start HTTP server", slog.String("address", a.server.Addr))
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	ctx, cancel := context.WithTimeout(context.Background(), config.App().HTTP.ShutdownTimeout())
	defer cancel()

	if err := a.closer.Close(ctx); err != nil {
		return fmt.Errorf("close deps: %w", err)
	}

	return nil
}
