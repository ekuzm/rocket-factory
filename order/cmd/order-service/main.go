package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/ekuzm/rocket-factory/order/internal/api/v1"
	"github.com/ekuzm/rocket-factory/order/internal/config"
	"github.com/ekuzm/rocket-factory/order/internal/integration/grpc/inventory"
	"github.com/ekuzm/rocket-factory/order/internal/integration/grpc/payment"
	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres"
	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres/transaction"
	"github.com/ekuzm/rocket-factory/order/internal/service"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	customMiddleware "github.com/ekuzm/rocket-factory/platform/pkg/middleware"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

func main() {
	if err := config.Setup(); err != nil {
		panic("Failed to setup config")
	}

	logger.Init(config.App().Logger.Level(), config.App().Logger.AsJSON())

	inventoryConn, err := grpc.NewClient(config.App().HTTP.InventoryAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service":           "order-service",
			"http_address":      config.App().HTTP.OrderAddress(),
			"inventory_address": config.App().HTTP.InventoryAddress(),
			"payment_address":   config.App().HTTP.PaymentAddress(),
			"component":         "inventory-client",
			"error":             err,
		}).Error("Failed to create client connection to inventory service")

		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			logger.WithFields(logrus.Fields{
				"service":           "order-service",
				"http_address":      config.App().HTTP.OrderAddress(),
				"inventory_address": config.App().HTTP.InventoryAddress(),
				"payment_address":   config.App().HTTP.PaymentAddress(),
				"component":         "inventory-client",
				"error":             err,
			}).Warn("Failed to close client connection to inventory service")
		}
	}()

	paymentConn, err := grpc.NewClient(config.App().HTTP.PaymentAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service":           "order-service",
			"http_address":      config.App().HTTP.OrderAddress(),
			"inventory_address": config.App().HTTP.InventoryAddress(),
			"payment_address":   config.App().HTTP.PaymentAddress(),
			"component":         "payment-client",
			"error":             err,
		}).Error("Failed to create client connection to payment service")

		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			logger.WithFields(logrus.Fields{
				"service":           "order-service",
				"http_address":      config.App().HTTP.OrderAddress(),
				"inventory_address": config.App().HTTP.InventoryAddress(),
				"payment_address":   config.App().HTTP.PaymentAddress(),
				"component":         "payment-client",
				"error":             err,
			}).Warn("Failed to close client connection to payment service")
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, config.App().Postgres.URI())
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service":           "order-service",
			"http_address":      config.App().HTTP.OrderAddress(),
			"inventory_address": config.App().HTTP.InventoryAddress(),
			"payment_address":   config.App().HTTP.PaymentAddress(),
			"component":         "postgres",
			"error":             err,
		}).Error("Failed to initialize pgx pool")

		return
	}
	defer func() {
		pool.Close()
		logger.WithFields(logrus.Fields{
			"service":           "order-service",
			"http_address":      config.App().HTTP.OrderAddress(),
			"inventory_address": config.App().HTTP.InventoryAddress(),
			"payment_address":   config.App().HTTP.PaymentAddress(),
			"component":         "postgres",
			"error":             err,
		}).Debug("Closed pgx pool")
	}()

	repository := postgres.New(ctx, pool)
	inventoryAdapter := inventory.New(inventoryV1.NewInventoryServiceClient(inventoryConn))
	paymentAdapter := payment.New(paymentV1.NewPaymentServiceClient(paymentConn))
	manager := transaction.New(pool)

	service := service.New(repository, inventoryAdapter, paymentAdapter, manager)
	api := api.New(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service":           "order-service",
			"http_address":      config.App().HTTP.OrderAddress(),
			"inventory_address": config.App().HTTP.InventoryAddress(),
			"payment_address":   config.App().HTTP.PaymentAddress(),
			"component":         "payment-client",
			"error":             err,
		}).Error("Failed to create order service server")

		return
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(customMiddleware.RequestLogger)

	router.Mount("/", orderServer)

	server := &http.Server{
		Addr:              config.App().HTTP.OrderAddress(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.WithFields(logrus.Fields{
			"service":           "order-service",
			"http_address":      config.App().HTTP.OrderAddress(),
			"inventory_address": config.App().HTTP.InventoryAddress(),
			"payment_address":   config.App().HTTP.PaymentAddress(),
		}).Debug("Starting order HTTP server...")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithFields(logrus.Fields{
				"service":           "order-service",
				"http_address":      config.App().HTTP.OrderAddress(),
				"inventory_address": config.App().HTTP.InventoryAddress(),
				"payment_address":   config.App().HTTP.PaymentAddress(),
				"error":             err,
			}).Error("Failed to listen and serve order HTTP server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	shutdownSignal := <-quit

	logger.WithFields(logrus.Fields{
		"service":           "order-service",
		"http_address":      config.App().HTTP.OrderAddress(),
		"inventory_address": config.App().HTTP.InventoryAddress(),
		"payment_address":   config.App().HTTP.PaymentAddress(),
		"signal":            shutdownSignal.String(),
		"error":             err,
	}).Warn("Shutting down the HTTP server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithFields(logrus.Fields{
			"service":           "order-service",
			"http_address":      config.App().HTTP.OrderAddress(),
			"inventory_address": config.App().HTTP.InventoryAddress(),
			"payment_address":   config.App().HTTP.PaymentAddress(),
			"error":             err,
		}).Error("Failed to shutdown the HTTP server")

		return
	}

	logger.WithFields(logrus.Fields{
		"service":           "order-service",
		"http_address":      config.App().HTTP.OrderAddress(),
		"inventory_address": config.App().HTTP.InventoryAddress(),
		"payment_address":   config.App().HTTP.PaymentAddress(),
	}).Debug("HTTP server successfully stopped")
}
