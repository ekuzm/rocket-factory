package model

type Payment struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethod
}
