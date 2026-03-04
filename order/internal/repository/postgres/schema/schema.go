package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type OrderRow struct {
	UUID            uuid.UUID      `db:"uuid"`
	UserUUID        uuid.UUID      `db:"user_uuid"`
	PartUUIDs       uuid.UUIDs     `db:"part_uuids"`
	TotalPrice      float64        `db:"total_price"`
	TransactionUUID sql.NullString `db:"transaction_uuid"`
	PaymentMethod   sql.NullString `db:"payment_method"`
	Status          Status         `db:"status"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       sql.NullTime   `db:"updated_at"`
}

func (or *OrderRow) Values() []any {
	return []any{or.UUID, or.UserUUID, or.PartUUIDs, or.TotalPrice, or.TransactionUUID, or.Status, or.PaymentMethod, or.CreatedAt, or.UpdatedAt}
}

const OrdersTable = "orders"

const (
	OrdersTableColumnUUID            = "uuid"
	OrdersTableColumnUserUUID        = "user_uuid"
	OrdersTableColumnPartUUIDs       = "part_uuids"
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
	OrdersTableColumnPartUUIDs,
	OrdersTableColumnTotalPrice,
	OrdersTableColumnTransactionUUID,
	OrdersTableColumnStatus,
	OrdersTableColumnPaymentMethod,
	OrdersTableColumnCreatedAt,
	OrdersTableColumnUpdatedAt,
}

type Status string

const (
	StatusPendingPayment Status = "pending_payment"
	StatusPaid           Status = "paid"
	StatusCancelled      Status = "cancelled"
)

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "unknown"
	PaymentMethodCard          PaymentMethod = "card"
	PaymentMethodSPB           PaymentMethod = "spb"
	PaymentMethodCreditCard    PaymentMethod = "credit_card"
	PaymentMethodInvestorMoney PaymentMethod = "investor_money"
)
