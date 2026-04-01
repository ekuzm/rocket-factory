package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ekuzm/rocket-factory/order/internal/config"
	"github.com/ekuzm/rocket-factory/platform/pkg/closer"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	customMiddleware "github.com/ekuzm/rocket-factory/platform/pkg/middleware"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

const (
	requestTimeout    = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 5 * time.Second
)

type app struct {
	di     *di
	server *http.Server
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
		slog.Debug("order service starting", "addr", config.App().HTTP.OrderAddress())

		errCh <- a.runHTTPServer()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
		defer shutdownCancel()

		if err := closer.CloseAll(shutdownCtx); err != nil {
			return err
		}

		slog.Debug("http server stopped")

		return nil
	case err := <-errCh:

		return err
	}
}

func (a *app) runHTTPServer() error {
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve order HTTP server: %w", err)
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

func (a *app) initServer(ctx context.Context) error {
	api, err := a.di.API(ctx)
	if err != nil {
		return fmt.Errorf("create order api: %w", err)
	}

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		return fmt.Errorf("create order server: %w", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(requestTimeout))
	router.Use(customMiddleware.RequestLogger)

	router.Mount("/", orderServer)

	a.server = &http.Server{
		Addr:              config.App().HTTP.OrderAddress(),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	closer.Add(a.server.Shutdown)

	return nil
}
