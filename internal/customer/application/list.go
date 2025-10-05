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
		Customers []CustomerView
	}

	ListCustomers struct {
		repo CustomerRepository
	}
)

func NewListCustomers(repo CustomerRepository) *ListCustomers {
	return &ListCustomers{repo: repo}
}

func (uc *ListCustomers) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
	items, err := uc.repo.List(ctx, in.Pager)
	if err != nil {
		return nil, err
	}

	customers := maps.Map(items, toView)
	return &ListOutput{Customers: customers}, nil
}
