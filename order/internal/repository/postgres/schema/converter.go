package entity

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/ekuzm/rocket-factory/order/internal/model"
)

var paymentMethodToSchema = map[model.PaymentMethod]PaymentMethod{
	model.PaymentMethodUnknown:       PaymentMethodUnknown,
	model.PaymentMethodCard:          PaymentMethodCard,
	model.PaymentMethodSPB:           PaymentMethodSPB,
	model.PaymentMethodCreditCard:    PaymentMethodCreditCard,
	model.PaymentMethodInvestorMoney: PaymentMethodInvestorMoney,
}

var statusToModel = map[Status]model.Status{
	StatusPendingPayment: model.StatusPendingPayment,
	StatusPaid:           model.StatusPaid,
	StatusCancelled:      model.StatusCancelled,
}

var paymentMethodToModel = map[PaymentMethod]model.PaymentMethod{
	PaymentMethodUnknown:       model.PaymentMethodUnknown,
	PaymentMethodCard:          model.PaymentMethodCard,
	PaymentMethodSPB:           model.PaymentMethodSPB,
	PaymentMethodCreditCard:    model.PaymentMethodCreditCard,
	PaymentMethodInvestorMoney: model.PaymentMethodInvestorMoney,
}

var statusToSchema = map[model.Status]Status{
	model.StatusPendingPayment: StatusPendingPayment,
	model.StatusPaid:           StatusPaid,
	model.StatusCancelled:      StatusCancelled,
}

func OrderToSchema(order model.Order) OrderRow {
	return OrderRow{
		UUID:            order.UUID,
		UserUUID:        order.Info.UserUUID,
		PartUUIDs:       order.Info.PartUUIDs,
		TotalPrice:      order.Info.TotalPrice,
		TransactionUUID: sql.NullString{String: order.Info.TransactionUUID.String(), Valid: order.Info.TransactionUUID != uuid.Nil},
		PaymentMethod:   sql.NullString{String: string(paymentMethodToSchema[order.Info.PaymentMethod]), Valid: order.Info.PaymentMethod != model.PaymentMethodUnknown},
		Status:          statusToSchema[order.Info.Status],
		CreatedAt:       order.CreatedAt,
	}
}

func OrderInfoToSchema(info model.OrderInfo) OrderRow {
	return OrderRow{
		UserUUID:        info.UserUUID,
		PartUUIDs:       info.PartUUIDs,
		TotalPrice:      info.TotalPrice,
		TransactionUUID: sql.NullString{String: info.TransactionUUID.String(), Valid: info.TransactionUUID != uuid.Nil},
		PaymentMethod:   sql.NullString{String: string(paymentMethodToSchema[info.PaymentMethod]), Valid: info.PaymentMethod != model.PaymentMethodUnknown},
		Status:          statusToSchema[info.Status],
	}
}

func OrderToModel(row OrderRow) (model.Order, error) {
	var transactionUUID uuid.UUID
	if row.TransactionUUID.Valid {
		uuid, err := uuid.Parse(row.TransactionUUID.String)
		if err != nil {
			return model.Order{}, fmt.Errorf("parse transaction UUID: %w", err)
		}

		transactionUUID = uuid
	}

	var updatedAt *time.Time
	if row.UpdatedAt.Valid {
		updatedAt = lo.ToPtr(row.UpdatedAt.Time)
	}

	return model.Order{
		UUID: row.UUID,
		Info: model.OrderInfo{
			UserUUID:        row.UserUUID,
			PartUUIDs:       row.PartUUIDs,
			TotalPrice:      row.TotalPrice,
			TransactionUUID: transactionUUID,
			PaymentMethod:   paymentMethodToModel[PaymentMethod(row.PaymentMethod.String)],
			Status:          statusToModel[row.Status],
		},
		CreatedAt: row.CreatedAt,
		UpdatedAt: updatedAt,
	}, nil
}
