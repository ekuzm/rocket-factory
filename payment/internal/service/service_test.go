package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	errs "github.com/ekuzm/rocket-factory/payment/internal/error"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
	"github.com/ekuzm/rocket-factory/payment/pkg/testutil"
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
			err: errs.ErrInvalidPaymentMethod,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := NewService()

			transactionUUID, err := service.PayOrder(test.args.ctx, test.args.orderUUID, test.args.userUUID, test.args.paymentMethod)
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
