package entity

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/ekuzm/rocket-factory/order/internal/model"
	"github.com/ekuzm/rocket-factory/shared/uuidx"
)

func OrderToSchema(order model.Order) OrderRow {
	return OrderRow{
		UUID:            order.UUID,
		UserUUID:        order.Info.UserUUID,
		TotalPrice:      order.Info.TotalPrice,
		TransactionUUID: sql.NullString{String: order.Info.TransactionUUID.String(), Valid: order.Info.TransactionUUID != uuid.Nil},
		PaymentMethod:   sql.NullInt32{Int32: int32(order.Info.PaymentMethod), Valid: order.Info.PaymentMethod != model.PaymentMethodUnknown},
		Status:          int(order.Info.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func OrderInfoToSchema(info model.OrderInfo) OrderRow {
	return OrderRow{
		UserUUID:        info.UserUUID,
		TotalPrice:      info.TotalPrice,
		TransactionUUID: sql.NullString{String: info.TransactionUUID.String(), Valid: info.TransactionUUID != uuid.Nil},
		PaymentMethod:   sql.NullInt32{Int32: int32(info.PaymentMethod), Valid: info.PaymentMethod != model.PaymentMethodUnknown},
		Status:          int(info.Status),
	}
}

func OrderToModel(row OrderRow, partUUIDs []string) (model.Order, error) {
	parsedPartUUIDs, err := uuidx.Parse(partUUIDs)
	if err != nil {
		return model.Order{}, fmt.Errorf("parse transaction UUIDs: %w", err)
	}

	var transactionUUID uuid.UUID
	if row.TransactionUUID.Valid {
		transactionUUID, err = uuid.Parse(row.TransactionUUID.String)
		if err != nil {
			return model.Order{}, fmt.Errorf("parse transaction UUID: %w", err)
		}
	}

	var updatedAt *time.Time
	if row.UpdatedAt.Valid {
		updatedAt = lo.ToPtr(row.UpdatedAt.Time)
	}

	return model.Order{
		UUID: row.UUID,
		Info: model.OrderInfo{
			UserUUID:        row.UserUUID,
			PartUUIDs:       parsedPartUUIDs,
			TotalPrice:      row.TotalPrice,
			TransactionUUID: transactionUUID,
			PaymentMethod:   model.PaymentMethod(int(row.PaymentMethod.Int32)),
			Status:          model.Status(row.Status),
		},
		CreatedAt: row.CreatedAt,
		UpdatedAt: updatedAt,
	}, nil
}
