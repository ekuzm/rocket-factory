package model

type Order struct {
	UUID            string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID string
	PaymentMethod   PaymentMethod
	Status          Status
}

type PaymentMethod int

const (
	PaymentMethodUnknown = iota
	PaymentMethodCard
	PaymentMethodSPB
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

type Status int

const (
	StatusPendingPayment = iota
	StatusPaid
	StatusCancelled
)
