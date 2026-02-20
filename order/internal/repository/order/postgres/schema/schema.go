package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type OrderRow struct {
	UUID            uuid.UUID      `db:"uuid"`
	UserUUID        uuid.UUID      `db:"user_uuid"`
	TotalPrice      float64        `db:"total_price"`
	TransactionUUID sql.NullString `db:"transaction_uuid"`
	PaymentMethod   sql.NullInt32  `db:"payment_method"`
	Status          int            `db:"status"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       sql.NullTime   `db:"updated_at"`
}

func (or *OrderRow) Values() []any {
	return []any{or.UUID, or.UserUUID, or.TotalPrice, or.TransactionUUID, or.Status, or.PaymentMethod, or.CreatedAt, or.UpdatedAt}
}

const OrdersTable = "orders"

const (
	OrdersTableColumnUUID            = "uuid"
	OrdersTableColumnUserUUID        = "user_uuid"
	OrdersTableColumnTotalPrice      = "total_price"
	OrdersTableColumnTransactionUUID = "transaction_uuid"
	OrdersTableColumnStatus          = "status"
	OrdersTableColumnPaymentMethod   = "payment_method"
	OrdersTableColumnCreatedAt       = "created_at"
	OrdersTableColumnUpdatedAt       = "updated_at"
)

var OrdersTableColumns = []string{
	OrdersTableColumnUUID,
	OrdersTableColumnUserUUID,
	OrdersTableColumnTotalPrice,
	OrdersTableColumnTransactionUUID,
	OrdersTableColumnStatus,
	OrdersTableColumnPaymentMethod,
	OrdersTableColumnCreatedAt,
	OrdersTableColumnUpdatedAt,
}

type OrderPartRow struct {
	OrderUUID string `db:"order_uuid"`
	PartUUID  string `db:"part_uuid"`
}

func (opr *OrderPartRow) Values() []any {
	return []any{opr.OrderUUID, opr.PartUUID}
}

const OrderPartsTable = "order_parts"

const (
	OrderPartsColumnOrderUUID = "order_uuid"
	OrderPartsColumnPartUUID  = "part_uuid"
)

var OrderPartsTableColumns = []string{
	OrderPartsColumnOrderUUID,
	OrderPartsColumnPartUUID,
}
