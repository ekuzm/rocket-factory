package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/inventory/internal/converter"
	"github.com/ekuzm/rocket-factory/inventory/internal/dto"
	"github.com/ekuzm/rocket-factory/inventory/internal/service/mock"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

func TestGetPart(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *inventoryV1.GetPartRequest
	}

	tests := []struct {
		name string
		args args
		want *inventoryV1.GetPartResponse
		err  error
		mock func(*mock.InventoryService, args)
	}{
		{
			name: "returns response from get part",
			args: args{
				ctx: context.Background(),
				req: &inventoryV1.GetPartRequest{
					Uuid: testPartUUID,
				},
			},
			want: &inventoryV1.GetPartResponse{
				Part: converter.PartToAPI(makePart(t)),
			},
			err: nil,
			mock: func(service *mock.InventoryService, args args) {
				input := dto.GetPartInput{
					UUID: args.req.Uuid,
				}

				output := dto.GetPartOutput{
					Part: makePart(t),
				}

				service.On("GetPart", args.ctx, input).Once().Return(output, nil)
			},
		},
		{
			name: "returns service error from get part",
			args: args{
				ctx: context.Background(),
				req: &inventoryV1.GetPartRequest{},
			},
			want: &inventoryV1.GetPartResponse{},
			err:  ErrService,
			mock: func(service *mock.InventoryService, args args) {
				input := dto.GetPartInput{
					UUID: args.req.Uuid,
				}

				var output dto.GetPartOutput

				service.On("GetPart", args.ctx, input).Once().Return(output, ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := mock.NewInventoryService(t)
			api := NewAPI(service)

			test.mock(service, test.args)

			resp, err := api.GetPart(test.args.ctx, test.args.req)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Nil(t, resp)

				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, resp)
		})
	}
}
