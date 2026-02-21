package v1_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	api "github.com/ekuzm/rocket-factory/inventory/internal/api/v1"
	"github.com/ekuzm/rocket-factory/inventory/internal/api/v1/dto"
	mockInventory "github.com/ekuzm/rocket-factory/inventory/internal/api/v1/mock"
	errs "github.com/ekuzm/rocket-factory/inventory/internal/error"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	"github.com/ekuzm/rocket-factory/inventory/pkg/testutil"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
)

func TestGetPart(t *testing.T) {
	type args struct {
		ctx context.Context
		req *inventoryV1.GetPartRequest
	}

	tests := []struct {
		name          string
		args          args
		want          *inventoryV1.GetPartResponse
		err           error
		mockInventory func(*mockInventory.InventoryService, args)
	}{
		{
			name: "returns response from get part",
			args: args{
				ctx: context.Background(),
				req: &inventoryV1.GetPartRequest{
					Uuid: testutil.TestPartUUID.String(),
				},
			},
			want: &inventoryV1.GetPartResponse{
				Part: dto.PartToAPI(testutil.MakePart(t)),
			},
			err: nil,
			mockInventory: func(service *mockInventory.InventoryService, args args) {
				part := testutil.MakePart(t)

				service.On("GetPart", args.ctx, testutil.TestPartUUID).Once().Return(part, nil)
			},
		},
		{
			name: "returns invalid uuid format error when part uuid is invalid",
			args: args{
				ctx: context.Background(),
				req: &inventoryV1.GetPartRequest{
					Uuid: uuid.Invalid.String(),
				},
			},
			want: nil,
			err:  errs.ErrInvalidUUIDFormat,
			mockInventory: func(service *mockInventory.InventoryService, args args) {
				service.AssertNotCalled(t, mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns service error from get part",
			args: args{
				ctx: context.Background(),
				req: &inventoryV1.GetPartRequest{
					Uuid: testutil.TestPartUUID.String(),
				},
			},
			want: nil,
			err:  testutil.ErrService,
			mockInventory: func(service *mockInventory.InventoryService, args args) {
				uuid, _ := uuid.Parse(args.req.Uuid)

				service.On("GetPart", args.ctx, uuid).Once().Return(model.Part{}, testutil.ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := mockInventory.NewInventoryService(t)
			api := api.New(service)

			test.mockInventory(service, test.args)

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

func TestListParts(t *testing.T) {
	type args struct {
		ctx context.Context
		req *inventoryV1.ListPartsRequest
	}

	tests := []struct {
		name          string
		args          args
		want          *inventoryV1.ListPartsResponse
		err           error
		mockInventory func(*mockInventory.InventoryService, args)
	}{
		{
			name: "returns response from list parts",
			args: args{
				ctx: context.Background(),
				req: &inventoryV1.ListPartsRequest{
					Filter: testutil.MakeFilterAPI(t),
				},
			},
			want: &inventoryV1.ListPartsResponse{
				Parts: dto.PartsToAPI([]model.Part{testutil.MakePart(t)}),
			},
			err: nil,
			mockInventory: func(service *mockInventory.InventoryService, args args) {
				filter := testutil.MakeFilter(t)
				parts := []model.Part{testutil.MakePart(t)}

				service.On("ListParts", args.ctx, filter).Once().Return(parts, nil)
			},
		},
		{
			name: "returns invalid uuid error when part uuids are invalid",
			args: args{
				ctx: context.Background(),
				req: &inventoryV1.ListPartsRequest{
					Filter: &inventoryV1.PartsFilter{
						Uuids: []string{uuid.Invalid.String()},
					},
				},
			},
			want: nil,
			err:  errs.ErrInvalidUUIDFormat,
			mockInventory: func(service *mockInventory.InventoryService, args args) {
				service.AssertNotCalled(t, "ListParts", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns service error from list parts",
			args: args{
				ctx: context.Background(),
				req: &inventoryV1.ListPartsRequest{
					Filter: testutil.MakeFilterAPI(t),
				},
			},
			want: nil,
			err:  testutil.ErrService,
			mockInventory: func(service *mockInventory.InventoryService, args args) {
				filter := testutil.MakeFilter(t)

				service.On("ListParts", args.ctx, filter).Once().Return(nil, testutil.ErrService)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := mockInventory.NewInventoryService(t)
			api := api.New(service)

			test.mockInventory(service, test.args)

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
