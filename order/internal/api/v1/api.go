package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ekuzm/rocket-factory/order/internal/api/v1/dto"
	errs "github.com/ekuzm/rocket-factory/order/internal/error"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	orderDto "github.com/ekuzm/rocket-factory/order/internal/service/dto"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
	"github.com/ekuzm/rocket-factory/shared/uuidx"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs uuid.UUIDs) (orderDto.Summary, error)
	GetOrder(ctx context.Context, uuid uuid.UUID) (model.Order, error)
	CancelOrder(ctx context.Context, uuid uuid.UUID) error
	PayOrder(ctx context.Context, uuid uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}

type api struct {
	service OrderService
	orderV1.UnimplementedHandler
}

func New(service OrderService) *api {
	return &api{
		service: service,
	}
}

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (*orderV1.CreateOrderResponse, error) {
	userUUID, err := uuid.Parse(req.UserUUID)
	if err != nil {
		return nil, fmt.Errorf("user UUID: %w", errs.ErrInvalidFormat)
	}
	partUUIDs, err := uuidx.Parse(req.PartUuids)
	if err != nil {
		return nil, fmt.Errorf("part UUIDs: %w", errs.ErrInvalidFormat)
	}

	summary, err := a.service.CreateOrder(ctx, userUUID, partUUIDs)
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	return &orderV1.CreateOrderResponse{UUID: summary.OrderUUID.String(), TotalPrice: summary.TotalPrice}, nil
}

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (*orderV1.Order, error) {
	uuid, err := uuid.Parse(params.OrderUUID)
	if err != nil {
		return nil, fmt.Errorf("order UUID: %w", errs.ErrInvalidFormat)
	}

	order, err := a.service.GetOrder(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	return dto.OrderToAPI(order), nil
}

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (*orderV1.NoContent, error) {
	uuid, err := uuid.Parse(params.OrderUUID)
	if err != nil {
		return nil, fmt.Errorf("order UUID: %w", errs.ErrInvalidFormat)
	}

	err = a.service.CancelOrder(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("cancel order: %w", err)
	}

	return &orderV1.NoContent{Code: http.StatusNoContent, Message: "order successfully cancelled"}, nil
}

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (*orderV1.PayOrderResponse, error) {
	uuid, err := uuid.Parse(params.OrderUUID)
	if err != nil {
		return nil, fmt.Errorf("order UUID: %w", errs.ErrInvalidFormat)
	}

	transactionUUID, err := a.service.PayOrder(ctx, uuid, dto.PaymentMethodToModel[req.PaymentMethod])
	if err != nil {
		return nil, fmt.Errorf("pay order: %w", err)
	}

	return &orderV1.PayOrderResponse{TransactionUUID: transactionUUID.String()}, nil
}

func (a *api) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	if err == nil {
		return nil
	}

	var (
		code    int
		message string
	)

	status, ok := status.FromError(err)
	if !ok {
		switch {
		case errors.Is(err, errs.ErrInvalidFormat):
			code, message = http.StatusBadRequest, err.Error()
		case errors.Is(err, errs.ErrOrderNotFound):
			code, message = http.StatusNotFound, err.Error()
		case errors.Is(err, errs.ErrStatusCancelled) || errors.Is(err, errs.ErrStatusPaid):
			code, message = http.StatusConflict, err.Error()
		default:
			code, message = http.StatusInternalServerError, err.Error()
		}
	} else {
		switch status.Code() {
		case codes.InvalidArgument:
			code, message = http.StatusBadRequest, status.Message()
		case codes.NotFound:
			code, message = http.StatusNotFound, status.Message()
		default:
			code, message = http.StatusInternalServerError, status.Message()
		}
	}

	return &orderV1.GenericErrorStatusCode{
		StatusCode: code,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(code),
			Message: orderV1.NewOptString(message),
		},
	}
}
