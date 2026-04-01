package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres/entity"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
)

type repository struct {
	pool Pool
}

func New(pool *pgxpool.Pool) *repository {
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
		slog.Error("order save failed", "orderUUID", order.UUID, "error", err)

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
			slog.Warn("order load missed", "orderUUID", uuid)

			return model.Order{}, fmt.Errorf("order: %w", errs.ErrNotFound)
		}

		slog.Error("order load failed", "orderUUID", uuid, "error", err)

		return model.Order{}, fmt.Errorf("execute select query from orders table: %w", err)
	}

	order, err := entity.OrderToModel(row)
	if err != nil {
		slog.Error("order convert failed", "orderUUID", uuid, "error", err)

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
		slog.Error("order update failed", "orderUUID", uuid, "error", err)

		return fmt.Errorf("execute update orders table: %w", err)
	}

	if res.RowsAffected() == 0 {
		slog.Warn("order update missed", "orderUUID", uuid)

		return fmt.Errorf("order: %w", errs.ErrNotFound)
	}

	return nil
}
