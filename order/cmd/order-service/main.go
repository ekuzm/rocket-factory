package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/ekuzm/rocket-factory/order/internal/api/v1"
	"github.com/ekuzm/rocket-factory/order/internal/integration/grpc/inventory"
	"github.com/ekuzm/rocket-factory/order/internal/integration/grpc/payment"
	customMiddleware "github.com/ekuzm/rocket-factory/order/internal/middleware"
	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres"
	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres/transaction"
	"github.com/ekuzm/rocket-factory/order/internal/service"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	inventoryServiceAddress = ":50051"
	paymentServiceAddress   = ":50052"
	orderServiceAddress     = ":8080"
	requestTimeout          = 10 * time.Second
	readHeaderTimeout       = 5 * time.Second
	shutdownTimeout         = 10 * time.Second
	dbTimeout               = 5 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(inventoryServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to get client connection to inventory service: %v", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("Failed to close client connection to inventory service: %v", cerr)
		}
	}()

	paymentConn, err := grpc.NewClient(paymentServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to get client connection to payment service: %v", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("Failed to close client connection to payment service: %v", cerr)
		}
	}()

	dbURI := os.Getenv("DB_URI")

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("Failed to initialize pgxpool: %v", err)
		return
	}
	defer pool.Close()

	repository := postgres.NewRepository(ctx, pool)
	inventoryAdapter := inventory.NewAdapter(inventoryV1.NewInventoryServiceClient(inventoryConn))
	paymentAdapter := payment.NewAdapter(paymentV1.NewPaymentServiceClient(paymentConn))
	manager := transaction.NewManager(pool)

	service := service.New(repository, inventoryAdapter, paymentAdapter, manager)
	api := api.New(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("Failed to create order service server: %v", err)
		return
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(requestTimeout))
	router.Use(customMiddleware.RequestLogger)

	router.Mount("/", orderServer)

	server := &http.Server{
		Addr:              orderServiceAddress,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("Order HTTP server listening at %s", orderServiceAddress)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to listen and serve order service HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Printf("Shutting down the HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown the HTTP server: %v", err)
		return
	}

	log.Printf("Successfully stopped HTTP server")
}
