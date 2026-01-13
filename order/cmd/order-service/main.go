package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

type OrderStorage struct {
	orders map[string]*orderV1.Order
	mtx    sync.RWMutex
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.Order),
		mtx:    sync.RWMutex{},
	}
}

func (o *OrderStorage) SaveOrder(order *orderV1.Order) {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	o.orders[order.UUID.String()] = order
}

func (o *OrderStorage) GetOrder(uuid string) *orderV1.Order {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	return o.orders[uuid]
}

func (o *OrderStorage) CancelOrder(order *orderV1.Order) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	order.Status = orderV1.OrderStatusCANCELLED
}

func (o *OrderStorage) PayOrder(order *orderV1.Order, transactionUUID uuid.UUID, paymentMethod orderV1.PaymentMethod) {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.NewOptUUID(transactionUUID)
	order.PaymentMethod = orderV1.NewOptPaymentMethod(paymentMethod)
}

type OrderService struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewOrderService(
	storage *OrderStorage,
	inventoryClient inventoryV1.InventoryServiceClient,
	paymentClient paymentV1.PaymentServiceClient,
) *OrderService {
	return &OrderService{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func (o *OrderService) ListParts(ctx context.Context, filter *inventoryV1.PartsFilter) ([]*inventoryV1.Part, error) {
	resp, err := o.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return resp.Parts, nil
}

func (o *OrderService) CalculateTotalPrice(ctx context.Context, parts []*inventoryV1.Part) float64 {
	var totalPrice float64

	for _, part := range parts {
		totalPrice += part.Price
	}

	return totalPrice
}

func (o *OrderService) SaveOrder(ctx context.Context, order *orderV1.Order) {
	o.storage.SaveOrder(order)
}

func (o *OrderService) GetOrder(ctx context.Context, uuid string) (*orderV1.Order, error) {
	order := o.storage.GetOrder(uuid)
	if order == nil {
		return nil, fmt.Errorf("the order with %s uuid not found", uuid)
	}

	return order, nil
}

func (o *OrderService) CancelOrder(ctx context.Context, order *orderV1.Order) error {
	if order.Status == orderV1.OrderStatusPAID {
		return fmt.Errorf("the order has payed and can't be canceled")
	}
	if order.Status == orderV1.OrderStatusCANCELLED {
		return fmt.Errorf("the order has already canceled")
	}

	o.storage.CancelOrder(order)

	return nil
}

func (o *OrderService) PayOrder(ctx context.Context, order *orderV1.Order, paymentMethod orderV1.PaymentMethod) (uuid.UUID, error) {
	resp, err := o.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		Uuid:     order.UUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: mapOrderToPaymentMethod(paymentMethod),
	})
	if err != nil {
		return uuid.Nil, err
	}

	if order.Status == orderV1.OrderStatusPAID {
		return uuid.Nil, fmt.Errorf("order has payed")
	}

	if order.Status == orderV1.OrderStatusCANCELLED {
		return uuid.Nil, fmt.Errorf("order has canceled")
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse transaction uuid: %v", err)
	}

	o.storage.PayOrder(order, transactionUUID, paymentMethod)

	return transactionUUID, nil
}

type OrderHandler struct {
	service *OrderService
}

func NewOrderHandler(service *OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (o *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	partsUUIDs := uuidsToStrings(req.GetPartUuids())

	filter := &inventoryV1.PartsFilter{
		Uuids: partsUUIDs,
	}

	parts, err := o.service.ListParts(ctx, filter)
	if err != nil {
		status := status.Convert(err)
		switch status.Code() {
		case codes.NotFound:
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "some parts not found:" + status.Message(),
			}, nil
		case codes.InvalidArgument:
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "bad request",
			}, nil
		default:
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "timeout when retrieving parts",
			}, nil
		}
	}

	totalPrice := o.service.CalculateTotalPrice(ctx, parts)
	orderUUID := uuid.New()

	order := &orderV1.Order{
		UUID:       orderUUID,
		UserUUID:   req.GetUserUUID(),
		PartUuids:  req.GetPartUuids(),
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	o.service.SaveOrder(ctx, order)

	return &orderV1.CreateOrderResponse{
		UUID:       orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

func (o *OrderHandler) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	if _, err := uuid.Parse(params.OrderUUID.String()); err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	order, err := o.service.GetOrder(ctx, params.OrderUUID.String())
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}, nil
	}

	return order, nil
}

func (o *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	if _, err := uuid.Parse(params.OrderUUID.String()); err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	order, err := o.service.GetOrder(ctx, params.OrderUUID.String())
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}, nil
	}

	if err = o.service.CancelOrder(ctx, order); err != nil {
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: err.Error(),
		}, nil
	}

	return &orderV1.NoContent{
		Code:    http.StatusNoContent,
		Message: fmt.Sprintf("order with %v uuid successfully canceled", params.OrderUUID),
	}, nil
}

func (o *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order, err := o.service.GetOrder(ctx, params.OrderUUID.String())
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}, nil
	}

	transactionUUID, err := o.service.PayOrder(
		ctx,
		order,
		req.GetPaymentMethod(),
	)
	if err != nil {
		status := status.Convert(err)
		switch status.Code() {
		case codes.InvalidArgument:
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		default:
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

func (o *OrderHandler) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func uuidsToStrings(uuids []uuid.UUID) []string {
	strings := make([]string, 0, len(uuids))

	for _, uuid := range uuids {
		strings = append(strings, uuid.String())
	}

	return strings
}

func mapOrderToPaymentMethod(orderPaymentMethod orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch orderPaymentMethod {
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SPB
	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

const (
	orderServiceHost        = "127.0.0.1"
	orderServicePort        = "8080"
	inventoryServiceAddress = "127.0.0.1:50051"
	paymentServiceAddress   = "127.0.0.1:50052"
	requestTimeout          = 10 * time.Second
	readHeaderTimeout       = 5 * time.Second
	shutdownTimeout         = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("Failed to close connection to inventory service: %v", cerr)
		}
	}()

	paymentConn, err := grpc.NewClient(
		paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to create payment service client: %v", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("Failed to close connection to payment service: %v", cerr)
		}
	}()

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	orderStorage := NewOrderStorage()
	orderService := NewOrderService(orderStorage, inventoryClient, paymentClient)
	orderHandler := NewOrderHandler(orderService)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("Failed to close connection to inventory service: %v", err)
		}
		log.Fatalf("Failed to intialize server: %v", err)
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(requestTimeout))

	router.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort(orderServiceHost, orderServicePort),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("Start listening http server on %s", net.JoinHostPort(orderServiceHost, orderServicePort))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to listening http server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown the server: %v", err)
	}

	log.Printf("Server stopped")
}
