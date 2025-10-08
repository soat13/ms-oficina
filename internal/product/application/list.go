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
	Products []ProductView
}

type ListProducts struct {
	repo Repository
}

func NewListProducts(repo Repository) *ListProducts {
	return &ListProducts{repo: repo}
}

func (uc *ListProducts) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
	items, err := uc.repo.List(ctx, in.Pager)
	if err != nil {
		return nil, err
	}

	return &ListOutput{
		Products: maps.Map(items, toView),
	}, nil
}
