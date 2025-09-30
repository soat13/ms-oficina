package application

import (
	"context"

	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/utils/pagination"
)

type ListInput struct {
	Limit  int
	Offset int
}

type ListOutput struct {
	Customers []CustomerView
}

type ListCustomers struct {
	repo CustomerRepository
}

func NewListCustomers(repo CustomerRepository) *ListCustomers {
	return &ListCustomers{repo: repo}
}

func (uc *ListCustomers) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
	limit, offset := pagination.LimitAndOffset(in.Limit, in.Offset)
	items, err := uc.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	customers := maps.Map(items, toView)
	return &ListOutput{Customers: customers}, nil
}
