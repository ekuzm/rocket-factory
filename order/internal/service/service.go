package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	api "github.com/ekuzm/rocket-factory/order/internal/api/v1"
	errs "github.com/ekuzm/rocket-factory/order/internal/error"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	"github.com/ekuzm/rocket-factory/order/internal/model/supplier"
	"github.com/ekuzm/rocket-factory/order/internal/service/dto"
)

type OrderRepository interface {
	Save(ctx context.Context, order model.Order) error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (model.Order, error)
	Update(ctx context.Context, uuid uuid.UUID, info model.OrderInfo) error
}

type InventoryPort interface {
	ListParts(ctx context.Context, filter supplier.Filter) ([]supplier.Part, error)
}

type PaymentPort interface {
	PayOrder(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}

type TransactionManager interface {
	Wrap(context.Context, func(context.Context) error) error
}

var _ api.OrderService = (*service)(nil)

type service struct {
	repository    OrderRepository
	inventoryPort InventoryPort
	paymentPort   PaymentPort
	manager       TransactionManager
}

func New(
	repository OrderRepository,
	inventoryPort InventoryPort,
	paymentPort PaymentPort,
	manager TransactionManager,
) *service {
	return &service{
		repository:    repository,
		inventoryPort: inventoryPort,
		paymentPort:   paymentPort,
		manager:       manager,
	}
}

func (s *service) CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs uuid.UUIDs) (dto.Summary, error) {
	parts, err := s.inventoryPort.ListParts(ctx, supplier.Filter{UUIDs: partUUIDs})
	if err != nil {
		return dto.Summary{}, fmt.Errorf("inventory port: %w", err)
	}

	var totalPrice float64

	for _, part := range parts {
		totalPrice += part.Price
	}

	order := model.Order{
		UUID: uuid.New(),
		Info: model.OrderInfo{
			UserUUID:   userUUID,
			PartUUIDs:  partUUIDs,
			TotalPrice: totalPrice,
			Status:     model.StatusPendingPayment,
		},
		CreatedAt: time.Now(),
	}

	err = s.manager.Wrap(ctx, func(ctx context.Context) error {
		if err := s.repository.Save(ctx, order); err != nil {
			return fmt.Errorf("repository: %w", err)
		}

		return nil
	})
	if err != nil {
		return dto.Summary{}, err
	}

	return dto.Summary{OrderUUID: order.UUID, TotalPrice: order.Info.TotalPrice}, nil
}

func (s *service) GetOrder(ctx context.Context, uuid uuid.UUID) (model.Order, error) {
	order, err := s.repository.GetByUUID(ctx, uuid)
	if err != nil {
		return model.Order{}, fmt.Errorf("repository: %w", err)
	}

	return order, nil
}

func (s *service) CancelOrder(ctx context.Context, uuid uuid.UUID) error {
	err := s.manager.Wrap(ctx, func(ctx context.Context) error {
		order, err := s.repository.GetByUUID(ctx, uuid)
		if err != nil {
			return fmt.Errorf("repository: %w", err)
		}

		if order.Info.Status == model.StatusCancelled {
			return errs.ErrStatusCancelled
		}

		if order.Info.Status == model.StatusPaid {
			return errs.ErrStatusPaid
		}

		order.Info.Status = model.StatusCancelled

		if err = s.repository.Update(ctx, uuid, order.Info); err != nil {
			return fmt.Errorf("repository: %w", err)
		}

		return nil
	})

	return err
}

func (s *service) PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	var transactionUUID uuid.UUID
	err := s.manager.Wrap(ctx, func(ctx context.Context) error {
		order, err := s.repository.GetByUUID(ctx, orderUUID)
		if err != nil {
			return fmt.Errorf("repository: %w", err)
		}

		if order.Info.Status == model.StatusCancelled {
			return errs.ErrStatusCancelled
		}

		if order.Info.Status == model.StatusPaid {
			return errs.ErrStatusPaid
		}

		transactionUUID, err = s.paymentPort.PayOrder(ctx, orderUUID, order.Info.UserUUID, paymentMethod)
		if err != nil {
			return fmt.Errorf("payment port: %w", err)
		}

		order.Info.PaymentMethod = paymentMethod
		order.Info.Status = model.StatusPaid
		order.Info.TransactionUUID = transactionUUID

		if err = s.repository.Update(ctx, orderUUID, order.Info); err != nil {
			return fmt.Errorf("repository: %w", err)
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return transactionUUID, nil
}
