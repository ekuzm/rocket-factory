package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres/entity"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
)

type repository struct {
	pool Pool
}

func New(ctx context.Context, pool *pgxpool.Pool) *repository {
	return &repository{
		pool: Pool{pool},
	}
}

func (r *repository) Save(ctx context.Context, order model.Order) error {
	row := entity.OrderToEntity(order)

	insertBuilder := sq.Insert(entity.OrdersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(entity.OrdersTableColumns...).
		Values(row.Values()...)

	if _, err := r.pool.Exec(ctx, insertBuilder); err != nil {
		logger.WithFields(logrus.Fields{
			"Order UUID":   order.UUID,
			"User UUID":    order.Info.UserUUID,
			"Part Count":   len(order.Info.PartUUIDs),
			"Order Status": order.Info.Status,
			"error":        err,
		}).Error("Failed to save order")

		return fmt.Errorf("execute insert query into orders table: %w", err)
	}

	return nil
}

func (r *repository) GetByUUID(ctx context.Context, uuid uuid.UUID) (model.Order, error) {
	selectBuilder := sq.Select(entity.OrdersTableColumns...).
		PlaceholderFormat(sq.Dollar).
		From(entity.OrdersTable).
		Where(sq.Eq{entity.OrdersTableColumnUUID: uuid})

	var row entity.OrderRow

	if err := r.pool.Get(ctx, &row, selectBuilder); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithFields(logrus.Fields{
				"Order UUID": uuid,
			}).Warn("Failed to get order by UUID, not found")

			return model.Order{}, fmt.Errorf("order: %w", errs.ErrNotFound)
		}

		logger.WithFields(logrus.Fields{
			"Order UUID": uuid,
			"error":      err,
		}).Error("Failed to get order by UUID")

		return model.Order{}, fmt.Errorf("execute select query from orders table: %w", err)
	}

	order, err := entity.OrderToModel(row)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Order UUID": uuid,
			"error":      err,
		}).Error("Failed to convert order entity to model")

		return model.Order{}, fmt.Errorf("entity to model: %w", err)
	}

	return order, nil
}

func (r *repository) Update(ctx context.Context, uuid uuid.UUID, info model.OrderInfo) error {
	row := entity.OrderInfoToEntity(info)

	updateBuilder := sq.Update(entity.OrdersTable).
		PlaceholderFormat(sq.Dollar).
		Set(entity.OrdersTableColumnUserUUID, row.UserUUID).
		Set(entity.OrdersTableColumnTotalPrice, row.TotalPrice).
		Set(entity.OrdersTableColumnTransactionUUID, row.TransactionUUID).
		Set(entity.OrdersTableColumnPaymentMethod, row.PaymentMethod).
		Set(entity.OrdersTableColumnStatus, row.Status).
		Set(entity.OrdersTableColumnUpdatedAt, sq.Expr("NOW()")).
		Where(sq.Eq{entity.OrdersTableColumnUUID: uuid})

	res, err := r.pool.Exec(ctx, updateBuilder)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Order UUID":       uuid,
			"User UUID":        info.UserUUID,
			"Order Status":     info.Status,
			"Payment Method":   info.PaymentMethod,
			"Transaction UUID": info.TransactionUUID,
			"error":            err,
		}).Error("Failed to update order")

		return fmt.Errorf("execute update orders table: %w", err)
	}

	if res.RowsAffected() == 0 {
		logger.WithFields(logrus.Fields{
			"Order UUID": uuid,
		}).Warn("Failed to update order, not found")

		return fmt.Errorf("order: %w", errs.ErrNotFound)
	}

	return nil
}
