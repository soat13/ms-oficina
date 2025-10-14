package application

import (
	"context"

	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/utils/pagination"
)

type (
	ListInput struct {
		Pager pagination.Pagination
	}

	ListOutput struct {
		Services []ServiceView
	}

	ListServices struct {
		repo ServiceRepository
	}
)

func NewListServices(repo ServiceRepository) *ListServices {
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
