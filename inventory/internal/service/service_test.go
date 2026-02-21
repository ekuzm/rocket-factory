package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	"github.com/ekuzm/rocket-factory/inventory/internal/service"
	mockInventory "github.com/ekuzm/rocket-factory/inventory/internal/service/mock"
	"github.com/ekuzm/rocket-factory/inventory/pkg/testutil"
)

func TestGetPart(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid uuid.UUID
	}

	tests := []struct {
		name string
		args args
		want model.Part
		err  error
		mock func(*mockInventory.InventoryRepository, args)
	}{
		{
			name: "returns part",
			args: args{
				ctx:  context.Background(),
				uuid: testutil.TestPartUUID,
			},
			want: testutil.MakePart(t),
			err:  nil,
			mock: func(repo *mockInventory.InventoryRepository, args args) {
				part := testutil.MakePart(t)

				repo.On("GetPart", args.ctx, args.uuid).Once().Return(part, nil)
			},
		},
		{
			name: "returns repository error from get part method",
			args: args{
				ctx:  context.Background(),
				uuid: testutil.TestPartUUID,
			},
			want: model.Part{},
			err:  testutil.ErrRepository,
			mock: func(repo *mockInventory.InventoryRepository, args args) {
				var part model.Part

				repo.On("GetPart", args.ctx, args.uuid).Once().Return(part, testutil.ErrRepository)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repo := mockInventory.NewInventoryRepository(t)
			service := service.New(repo)

			test.mock(repo, test.args)

			output, err := service.GetPart(test.args.ctx, test.args.uuid)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Empty(t, output)

				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, output)
		})
	}
}

func TestListParts(t *testing.T) {
	type args struct {
		ctx    context.Context
		filter model.Filter
	}

	tests := []struct {
		name string
		args args
		want []model.Part
		err  error
		mock func(*mockInventory.InventoryRepository, args)
	}{
		{
			name: "returns parts",
			args: args{
				ctx:    context.Background(),
				filter: testutil.MakeFilter(t),
			},
			want: []model.Part{testutil.MakePart(t)},
			err:  nil,
			mock: func(repo *mockInventory.InventoryRepository, args args) {
				parts := []model.Part{testutil.MakePart(t)}

				repo.On("ListParts", args.ctx, args.filter).Once().Return(parts, nil)
			},
		},
		{
			name: "returns repository error from list parts method",
			args: args{
				ctx:    context.Background(),
				filter: testutil.MakeFilter(t),
			},
			want: nil,
			err:  testutil.ErrRepository,
			mock: func(repo *mockInventory.InventoryRepository, args args) {
				repo.On("ListParts", args.ctx, args.filter).Once().Return(nil, testutil.ErrRepository)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repo := mockInventory.NewInventoryRepository(t)
			service := service.New(repo)

			test.mock(repo, test.args)

			output, err := service.ListParts(test.args.ctx, test.args.filter)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)
				require.Empty(t, output)

				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, output)
		})
	}
}
