package application

import (
	"context"

	"github.com/soat13/fase-1-oficina/pkg/maps"
)

type ListInput struct {
	Limit  int
	Offset int
}

type ListOutput struct {
	Customers []CustomerView `json:"customers"`
}

type ListCustomers struct {
	repo CustomerRepository
}

func NewListCustomers(repo CustomerRepository) *ListCustomers {
	return &ListCustomers{repo: repo}
}

func (uc *ListCustomers) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
	limit, offset := limitAndOffset(in)
	items, err := uc.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	customers := maps.Map(items, toView)
	return &ListOutput{Customers: customers}, nil
}

func limitAndOffset(input ListInput) (int, int) {
	if input.Limit <= 0 {
		input.Limit = 50
	}
	if input.Offset < 0 {
		input.Offset = 0
	}

	return input.Limit, input.Offset
}
