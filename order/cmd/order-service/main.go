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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/ekuzm/rocket-factory/order/internal/api/order/v1"
	"github.com/ekuzm/rocket-factory/order/internal/client/grpc/inventory"
	"github.com/ekuzm/rocket-factory/order/internal/client/grpc/payment"
	customMiddleware "github.com/ekuzm/rocket-factory/order/internal/middleware"
	orderRepository "github.com/ekuzm/rocket-factory/order/internal/repository/order/memory"
	orderService "github.com/ekuzm/rocket-factory/order/internal/service/order"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
)

const (
	InventoryServiceAddress = "127.0.0.1:50051"
	PaymentServiceAddress   = "127.0.0.1:50052"
	OrderServiceAddress     = "127.0.0.1:8080"
	RequestTimeout          = 10 * time.Second
	ReadHeaderTimeout       = 5 * time.Second
	ShutdownTimeout         = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(InventoryServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to get client connection to inventory service: %v", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("Failed to close client connection to inventory service: %v", cerr)
		}
	}()

	paymentConn, err := grpc.NewClient(PaymentServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to get client connection to payment service: %v", err)
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("Failed to close client connection to payment service: %v", cerr)
		}
	}()

	repository := orderRepository.NewRepository()
	inventoryClient := inventory.NewClient(inventoryConn)
	paymentClient := payment.NewClient(paymentConn)

	service := orderService.NewService(repository, inventoryClient, paymentClient)
	api := orderAPI.NewAPI(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("Failed to create order service server: %v", err)
		return
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(RequestTimeout))
	router.Use(customMiddleware.RequestLogger)

	router.Mount("/", orderServer)

	server := &http.Server{
		Addr:              OrderServiceAddress,
		Handler:           router,
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

	go func() {
		log.Printf("Order HTTP server listening at %s", OrderServiceAddress)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to listen and serve order service HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Printf("Shutting down the HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown the HTTP server: %v", err)
		return
	}

	log.Printf("Successfully stopped HTTP server")
}
