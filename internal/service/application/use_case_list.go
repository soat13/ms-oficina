package application

import (
	"context"

	"github.com/soat13/oficina-utils/pkg/maps"
	"github.com/soat13/oficina-utils/pkg/pagination"
)

type (
	ListInput struct {
		Pager pagination.Pagination
	}

	ListOutput struct {
		Services []ServiceView
	}

	ListServices struct {
		repo Repository
	}
)

func NewListServices(repo Repository) *ListServices {
	return &ListServices{repo: repo}
}

func (uc *ListServices) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
	items, err := uc.repo.List(ctx, in.Pager)
	if err != nil {
		return nil, err
	}

	services := maps.Map(items, toView)
	return &ListOutput{Services: services}, nil
}
