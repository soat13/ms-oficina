package application

import (
	"context"

	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/utils/pagination"
)

type ListInput struct {
	Pager pagination.Pagination
}

type ListOutput struct {
	Services []ServiceView
}

type ListServices struct {
	repo Repository
}

func NewListServices(repo Repository) *ListServices {
	return &ListServices{repo: repo}
}

func (uc *ListServices) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	items, err := uc.repo.List(ctx, input.Pager)
	if err != nil {
		return nil, err
	}

	return &ListOutput{
		Services: maps.Map(items, toView),
	}, nil
}
