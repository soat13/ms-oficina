package application

import "context"

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
	if in.Limit <= 0 {
		in.Limit = 50
	}
	if in.Offset < 0 {
		in.Offset = 0
	}
	items, err := uc.repo.List(ctx, in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}
	out := make([]CustomerView, 0, len(items))
	for _, c := range items {
		out = append(out, toView(c))
	}
	return &ListOutput{Customers: out}, nil
}
