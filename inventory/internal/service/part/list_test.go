package part

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ekuzm/rocket-factory/inventory/internal/converter"
	"github.com/ekuzm/rocket-factory/inventory/internal/dto"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	repoConverter "github.com/ekuzm/rocket-factory/inventory/internal/repository/converter"
	repositoryMock "github.com/ekuzm/rocket-factory/inventory/internal/repository/mock"
)

func TestListParts(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.ListPartsInput
	}

	tests := []struct {
		name string
		args args
		want dto.ListPartsOutput
		err  error
		mock func(*repositoryMock.InventoryRepository, args)
	}{
		{
			name: "returns parts",
			args: args{
				ctx: context.Background(),
				input: dto.ListPartsInput{
					Filter: converter.FilterToModel(makeFilter(t, testPartUUID)),
				},
			},
			want: dto.ListPartsOutput{
				Parts: []model.Part{makePart(t)},
			},
			err: nil,
			mock: func(repo *repositoryMock.InventoryRepository, args args) {
				parts := []model.Part{makePart(t)}

				repo.On("ListParts", args.ctx, repoConverter.FilterToRepoModel(args.input.Filter)).Once().Return(parts, nil)
			},
		},
		{
			name: "returns invalid format error when filter is empty",
			args: args{
				ctx:   context.Background(),
				input: dto.ListPartsInput{},
			},
			want: dto.ListPartsOutput{},
			err:  model.ErrInvalidFormat,
			mock: func(repo *repositoryMock.InventoryRepository, args args) {
				repo.AssertNotCalled(t, "ListPart", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns invalid format error when filter has invalid uuid",
			args: args{
				ctx: context.Background(),
				input: dto.ListPartsInput{
					Filter: converter.FilterToModel(makeFilter(t, uuid.Invalid.String())),
				},
			},
			want: dto.ListPartsOutput{},
			err:  model.ErrInvalidFormat,
			mock: func(repo *repositoryMock.InventoryRepository, args args) {
				repo.AssertNotCalled(t, "ListPart", mock.Anything, mock.Anything)
			},
		},
		{
			name: "returns repository error from list parts method",
			args: args{
				ctx: context.Background(),
				input: dto.ListPartsInput{
					Filter: converter.FilterToModel(makeFilter(t, testPartUUID)),
				},
			},
			want: dto.ListPartsOutput{},
			err:  ErrRepository,
			mock: func(repo *repositoryMock.InventoryRepository, args args) {
				repo.On("ListParts", args.ctx, repoConverter.FilterToRepoModel(args.input.Filter)).Once().Return(nil, ErrRepository)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repo := repositoryMock.NewInventoryRepository(t)
			service := NewService(repo)

			test.mock(repo, test.args)

			output, err := service.ListParts(test.args.ctx, test.args.input)
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
