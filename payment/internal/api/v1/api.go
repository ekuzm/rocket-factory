package v1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/ekuzm/rocket-factory/payment/internal/api/v1/dto"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

type PaymentService interface {
	PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}

type api struct {
	paymentV1.UnimplementedPaymentServiceServer
	service PaymentService
}

func New(service PaymentService) *api {
	return &api{
		service: service,
	}
}

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	orderUUID, err := uuid.Parse(req.Uuid)
	if err != nil {
		slog.Warn("order uuid parse failed", "orderUUID", req.Uuid, "error", err)

		return nil, fmt.Errorf("parse order UUID: %w", errs.ErrInvalid)
	}

	transactionUUID, err := a.service.PayOrder(ctx, orderUUID, dto.PaymentMethodToModel[req.PaymentMethod])
	if err != nil {
		return nil, fmt.Errorf("payment service: %w", err)
	}

	return &paymentV1.PayOrderResponse{TransactionUuid: transactionUUID.String()}, nil
}
