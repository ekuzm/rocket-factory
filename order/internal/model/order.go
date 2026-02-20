package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	UUID      uuid.UUID
	Info      OrderInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type OrderInfo struct {
	UserUUID        uuid.UUID
	PartUUIDs       uuid.UUIDs
	TotalPrice      float64
	TransactionUUID uuid.UUID
	PaymentMethod   PaymentMethod
	Status          Status
}

type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSPB
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

type Status int

const (
	StatusPendingPayment Status = iota
	StatusPaid
	StatusCancelled
)
