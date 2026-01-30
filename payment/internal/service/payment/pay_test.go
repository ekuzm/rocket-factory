package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/payment/internal/dto"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx   context.Context
		input dto.PayOrderInput
	}

	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "returns transaction UUID",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     testOrderUUID,
					UserUUID:      testUserUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: nil,
		},
		{
			name: "returns invalid format error when order uuid is invalid",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     uuid.Invalid.String(),
					UserUUID:      testUserUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: model.ErrInvalidFormat,
		},
		{
			name: "returns invalid format error when order uuid is empty",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					UserUUID:      testUserUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: model.ErrInvalidFormat,
		},
		{
			name: "returns invalid format error when user uuid is invalid",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     testOrderUUID,
					UserUUID:      uuid.Invalid.String(),
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: model.ErrInvalidFormat,
		},
		{
			name: "returns invalid format error when user uuid is empty",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     testOrderUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			err: model.ErrInvalidFormat,
		},
		{
			name: "returns invalid format error when payment method is unknown",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID:     testOrderUUID,
					UserUUID:      testUserUUID,
					PaymentMethod: model.PaymentMethodUnknown,
				},
			},
			err: model.ErrInvalidFormat,
		},
		{
			name: "returns invalid format error when payment method is empty",
			args: args{
				ctx: context.Background(),
				input: dto.PayOrderInput{
					OrderUUID: testOrderUUID,
					UserUUID:  testUserUUID,
				},
			},
			err: model.ErrInvalidFormat,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewService()

			output, err := service.PayOrder(test.args.ctx, test.args.input)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Empty(t, output)

				return
			}

			require.NoError(t, err)
			require.NoError(t, uuid.Validate(output.TransactionUUID))
		})
	}
}
