package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/payment/internal/model"
	"github.com/ekuzm/rocket-factory/payment/internal/service"
	"github.com/ekuzm/rocket-factory/payment/pkg/testutil"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
)

func TestPayOrder(t *testing.T) {
	type args struct {
		ctx           context.Context
		orderUUID     uuid.UUID
		userUUID      uuid.UUID
		paymentMethod model.PaymentMethod
	}

	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "returns transaction UUID",
			args: args{
				ctx:           context.Background(),
				orderUUID:     testutil.TestOrderUUID,
				userUUID:      testutil.TestUserUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			err: nil,
		},
		{
			name: "returns invalid payment method error when payment method is an unknown",
			args: args{
				ctx:           context.Background(),
				orderUUID:     testutil.TestOrderUUID,
				userUUID:      testutil.TestUserUUID,
				paymentMethod: model.PaymentMethodUnknown,
			},
			err: errs.ErrInvalid,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := service.New()

			transactionUUID, err := service.PayOrder(test.args.ctx, test.args.orderUUID, test.args.paymentMethod)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Empty(t, transactionUUID)

				return
			}

			require.NoError(t, err)
			require.NoError(t, uuid.Validate(transactionUUID.String()))
		})
	}
}
