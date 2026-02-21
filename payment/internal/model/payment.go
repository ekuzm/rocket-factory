package model

type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSPB
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)
