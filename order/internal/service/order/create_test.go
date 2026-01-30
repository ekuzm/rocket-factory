package order

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	clientMock "github.com/ekuzm/rocket-factory/order/internal/client/mock"
	"github.com/ekuzm/rocket-factory/order/internal/dto"
	"github.com/ekuzm/rocket-factory/order/internal/model"
	repoMock "github.com/ekuzm/rocket-factory/order/internal/repository/mock"
	repoModel "github.com/ekuzm/rocket-factory/order/internal/repository/model"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx   context.Context
		input dto.CreateOrderInput
	}

	tests := []struct {
		name string
		args args
		want dto.CreateOrderOutput
		err  error
		mock func(*repoMock.OrderRepository, *clientMock.InventoryClient, args)
	}{
		{
			name: "returns UUID and total price",
			args: args{
				ctx: context.Background(),
				input: dto.CreateOrderInput{
					UserUUID: testUserUUID,
					PartUUIDs: []string{
						testPartUUID,
					},
				},
			},
			want: dto.CreateOrderOutput{
				TotalPrice: testTotalPrice,
			},
			err: nil,
			mock: func(repo *repoMock.OrderRepository, client *clientMock.InventoryClient, args args) {
				filter := model.Filter{
					UUIDs: args.input.PartUUIDs,
				}

				parts := []model.Part{makePart(t)}

				total := testTotalPrice

				client.On("ListParts", args.ctx, filter).Once().Return(parts, nil)
				repo.On("SaveOrder", args.ctx, mock.MatchedBy(func(order repoModel.Order) bool {
					return matchOrderInput(t, order, args.input, total)
				})).Once().Return(nil)
			},
		},
		{
			name: "returns inventory client error",
			args: args{
				ctx: context.Background(),
				input: dto.CreateOrderInput{
					UserUUID: testUserUUID,
					PartUUIDs: []string{
						testPartUUID,
					},
				},
			},
			want: dto.CreateOrderOutput{},
			err:  ErrInventoryClient,
			mock: func(repo *repoMock.OrderRepository, client *clientMock.InventoryClient, args args) {
				filter := model.Filter{
					UUIDs: args.input.PartUUIDs,
				}

				client.On("ListParts", args.ctx, filter).Once().Return(nil, ErrInventoryClient)
				repo.AssertNotCalled(t, "SaveOrder", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns repository error from get order method",
			args: args{
				ctx: context.Background(),
				input: dto.CreateOrderInput{
					UserUUID: testUserUUID,
					PartUUIDs: []string{
						testPartUUID,
					},
				},
			},
			want: dto.CreateOrderOutput{},
			err:  ErrRepository,
			mock: func(repo *repoMock.OrderRepository, client *clientMock.InventoryClient, args args) {
				filter := model.Filter{
					UUIDs: args.input.PartUUIDs,
				}

				parts := []model.Part{makePart(t)}

				client.On("ListParts", args.ctx, filter).Once().Return(parts, nil)
				repo.On("SaveOrder", args.ctx, mock.MatchedBy(func(order repoModel.Order) bool {
					return matchOrderInput(t, order, args.input, testPartPrice)
				})).Once().Return(ErrRepository)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repoMock.NewOrderRepository(t)
			inventoryClient := clientMock.NewInventoryClient(t)
			paymentClient := clientMock.NewPaymentClient(t)

			service := NewService(repo, inventoryClient, paymentClient)

			test.mock(repo, inventoryClient, test.args)

			output, err := service.CreateOrder(test.args.ctx, test.args.input)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Empty(t, output)

				return
			}

			require.NoError(t, err)
			require.NoError(t, uuid.Validate(output.OrderUUID))
			require.Equal(t, test.want.TotalPrice, output.TotalPrice)
		})
	}
}
