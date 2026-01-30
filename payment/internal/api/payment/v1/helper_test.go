package v1

import "errors"

const (
	testOrderUUID       = "11111111-1111-1111-1111-111111111111"
	testUserUUID        = "22222222-2222-2222-2222-222222222222"
	testTransactionUUID = "33333333-3333-3333-3333-333333333333"
)

var ErrService = errors.New("payment service error")
