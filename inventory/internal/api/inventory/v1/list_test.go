package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/inventory/internal/converter"
	"github.com/ekuzm/rocket-factory/inventory/internal/dto"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	"github.com/ekuzm/rocket-factory/inventory/internal/service/mock"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

func TestListParts(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *inventoryV1.ListPartsRequest
	}

	tests := []struct {
		name string
		args args
		want *inventoryV1.ListPartsResponse
		err  error
		mock func(*mock.InventoryService, args)
	}{
		{
			name: "returns response from list parts",
			args: args{
				ctx: context.Background(),
				req: &inventoryV1.ListPartsRequest{
					Filter: makeFilter(t, testPartUUID),
				},
			},
			want: &inventoryV1.ListPartsResponse{
				Parts: []*inventoryV1.Part{converter.PartToAPI(makePart(t))},
			},
			err: nil,
			mock: func(service *mock.InventoryService, args args) {
				input := dto.ListPartsInput{
					Filter: converter.FilterToModel(makeFilter(t, testPartUUID)),
				}

				output := dto.ListPartsOutput{
					Parts: []model.Part{makePart(t)},
				}

				service.On("ListParts", args.ctx, input).Once().Return(output, nil)
			},
		},
		{
			name: "returns service error from list parts",
			args: args{
				ctx: context.Background(),
				req: &inventoryV1.ListPartsRequest{},
			},
			want: nil,
			err:  ErrService,
			mock: func(service *mock.InventoryService, args args) {
				input := dto.ListPartsInput{}

				output := dto.ListPartsOutput{
					Parts: []model.Part{makePart(t)},
				}

				service.On("ListParts", args.ctx, input).Once().Return(output, ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := mock.NewInventoryService(t)
			api := NewAPI(service)

			test.mock(service, test.args)

			resp, err := api.ListParts(test.args.ctx, test.args.req)
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
