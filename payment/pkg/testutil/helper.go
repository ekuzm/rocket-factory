package testutil

import (
	"errors"

	"github.com/google/uuid"
)

var (
	TestOrderUUID       = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	TestUserUUID        = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	TestTransactionUUID = uuid.MustParse("33333333-3333-3333-3333-333333333333")
)

var ErrService = errors.New("payment service error")
