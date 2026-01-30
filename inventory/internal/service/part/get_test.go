package part

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/inventory/internal/dto"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	repositoryMock "github.com/ekuzm/rocket-factory/inventory/internal/repository/mock"
)

func TestGetPart(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx   context.Context
		input dto.GetPartInput
	}

	tests := []struct {
		name string
		args args
		want dto.GetPartOutput
		err  error
		mock func(*repositoryMock.InventoryRepository, args)
	}{
		{
			name: "returns part",
			args: args{
				ctx: context.Background(),
				input: dto.GetPartInput{
					UUID: testPartUUID,
				},
			},
			want: dto.GetPartOutput{
				Part: makePart(t),
			},
			err: nil,
			mock: func(repo *repositoryMock.InventoryRepository, args args) {
				part := makePart(t)

				repo.On("GetPart", args.ctx, args.input.UUID).Once().Return(part, nil)
			},
		},
		{
			name: "returns invalid format error when part uuid is empty",
			args: args{
				ctx:   context.Background(),
				input: dto.GetPartInput{},
			},
			want: dto.GetPartOutput{},
			err:  model.ErrInvalidFormat,
			mock: func(repo *repositoryMock.InventoryRepository, args args) {
				repo.AssertNotCalled(t, "GetPart", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns invalid format error when part uuid is invalid",
			args: args{
				ctx: context.Background(),
				input: dto.GetPartInput{
					UUID: uuid.Invalid.String(),
				},
			},
			want: dto.GetPartOutput{},
			err:  model.ErrInvalidFormat,
			mock: func(repo *repositoryMock.InventoryRepository, args args) {
				repo.AssertNotCalled(t, "GetPart", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns repository error from get part method",
			args: args{
				ctx: context.Background(),
				input: dto.GetPartInput{
					UUID: testPartUUID,
				},
			},
			want: dto.GetPartOutput{},
			err:  ErrRepository,
			mock: func(repo *repositoryMock.InventoryRepository, args args) {
				var part model.Part

				repo.On("GetPart", args.ctx, args.input.UUID).Once().Return(part, ErrRepository)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repositoryMock.NewInventoryRepository(t)
			service := NewService(repo)

			test.mock(repo, test.args)

			output, err := service.GetPart(test.args.ctx, test.args.input)
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
