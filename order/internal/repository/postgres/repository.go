package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	errs "github.com/ekuzm/rocket-factory/order/internal/error"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	schema "github.com/ekuzm/rocket-factory/order/internal/repository/postgres/schema"
)

type repository struct {
	pool Pool
}

func New(ctx context.Context, pool *pgxpool.Pool) *repository {
	return &repository{
		pool: Pool{pool},
	}
}

func (r *repository) SaveOrder(ctx context.Context, order model.Order) error {
	row := schema.OrderToSchema(order)

	insertBuilder := sq.Insert(schema.OrdersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(schema.OrdersTableColumns...).
		Values(row.Values()...)

	if _, err := r.pool.Exec(ctx, insertBuilder); err != nil {
		return fmt.Errorf("execute insert query into orders table: %w", err)
	}

	insertBuilder = sq.Insert(schema.OrderPartsTable).
		PlaceholderFormat(sq.Dollar).
		Columns(schema.OrderPartsTableColumns...)

	for _, partUUID := range order.Info.PartUUIDs {
		insertBuilder = insertBuilder.Values(order.UUID, partUUID)
	}

	if _, err := r.pool.Exec(ctx, insertBuilder); err != nil {
		return fmt.Errorf("execute insert query into order parts table: %w", err)
	}

	return nil
}

func (r *repository) GetOrder(ctx context.Context, uuid uuid.UUID) (model.Order, error) {
	selectBuilder := sq.Select(schema.OrdersTableColumns...).
		PlaceholderFormat(sq.Dollar).
		From(schema.OrdersTable).
		Where(sq.Eq{schema.OrdersTableColumnUUID: uuid})

	var row schema.OrderRow

	if err := r.pool.Get(ctx, &row, selectBuilder); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, errs.ErrOrderNotFound
		}

		return model.Order{}, fmt.Errorf("execute select query from orders table: %w", err)
	}

	selectBuilder = sq.Select(schema.OrderPartsColumnPartUUID).
		PlaceholderFormat(sq.Dollar).
		From(schema.OrderPartsTable).
		Where(sq.Eq{schema.OrderPartsColumnOrderUUID: uuid})

	var partUUIDs []string

	if err := r.pool.Select(ctx, &partUUIDs, selectBuilder); err != nil {
		return model.Order{}, fmt.Errorf("execute select query from order parts table: %w", err)
	}

	return schema.OrderToModel(row, partUUIDs)
}

func (r *repository) UpdateOrder(ctx context.Context, uuid uuid.UUID, info model.OrderInfo) error {
	row := schema.OrderInfoToSchema(info)

	updateBuilder := sq.Update(schema.OrdersTable).
		PlaceholderFormat(sq.Dollar).
		Set(schema.OrdersTableColumnUserUUID, row.UserUUID).
		Set(schema.OrdersTableColumnTotalPrice, row.TotalPrice).
		Set(schema.OrdersTableColumnTransactionUUID, row.TransactionUUID).
		Set(schema.OrdersTableColumnPaymentMethod, row.PaymentMethod).
		Set(schema.OrdersTableColumnStatus, row.Status).
		Set(schema.OrdersTableColumnUpdatedAt, sq.Expr("NOW()")).
		Where(sq.Eq{schema.OrdersTableColumnUUID: uuid})

	res, err := r.pool.Exec(ctx, updateBuilder)
	if err != nil {
		return fmt.Errorf("execute update orders table: %w", err)
	}

	if res.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}

	deleteBuilder := sq.Delete(schema.OrderPartsTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{schema.OrderPartsColumnOrderUUID: uuid})

	if _, err = r.pool.Exec(ctx, deleteBuilder); err != nil {
		return fmt.Errorf("execute delete columns from order parts table: %w", err)
	}

	insertBuilder := sq.Insert(schema.OrderPartsTable).
		PlaceholderFormat(sq.Dollar).
		Columns(schema.OrderPartsTableColumns...)

	if len(info.PartUUIDs) > 0 {
		for _, partUUID := range info.PartUUIDs {
			insertBuilder = insertBuilder.
				Values(uuid, partUUID)
		}

		if _, err = r.pool.Exec(ctx, insertBuilder); err != nil {
			return fmt.Errorf("execute insert into order parts table: %w", err)
		}
	}

	return nil
}
