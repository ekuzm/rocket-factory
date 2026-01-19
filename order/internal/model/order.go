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

type PaymentMethod string

const (
	PaymentMethodUnknown       = "UNKNOWN"
	PaymentMethodCard          = "CARD"
	PaymentMethodSPB           = "SPB"
	PaymentMethodCreditCard    = "CREDIT_CARD"
	PaymentMethodInvestorMoney = "INVESTOR_MONEY"
)

type Status string

const (
	StatusPendingPayment = "PENDING_PAYMENT"
	StatusPaid           = "PAID"
	StatusCancelled      = "CANCELLED"
)
