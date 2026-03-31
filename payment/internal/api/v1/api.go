package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ekuzm/rocket-factory/payment/internal/api/v1/dto"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

type PaymentService interface {
	PayOrder(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
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
		logger.WithFields(logrus.Fields{
			"Order UUID": req.Uuid,
			"error":      err,
		}).Warn("Failed to parse order UUID")

		return nil, fmt.Errorf("parse order UUID: %w", errs.ErrInvalid)
	}
	userUUID, err := uuid.Parse(req.UserUuid)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Order UUID": req.Uuid,
			"User UUID":  req.UserUuid,
			"error":      err,
		}).Warn("Failed to parse user UUID")

		return nil, fmt.Errorf("parse user UUID: %w", errs.ErrInvalid)
	}

	transactionUUID, err := a.service.PayOrder(ctx, orderUUID, userUUID, dto.PaymentMethodToModel[req.PaymentMethod])
	if err != nil {
		return nil, fmt.Errorf("payment service: %w", err)
	}

	return &paymentV1.PayOrderResponse{TransactionUuid: transactionUUID.String()}, nil
}
