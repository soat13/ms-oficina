package application

import "context"

type ListInput struct {
	Limit  int
	Offset int
}

type ListOutput struct {
	Products []ProductView `json:"products"`
}

type ListProducts struct {
	repo ProductRepository
}

func NewListProducts(repo ProductRepository) *ListProducts {
	return &ListProducts{repo: repo}
}

func (uc *ListProducts) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
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
	out := make([]ProductView, 0, len(items))
	for _, p := range items {
		out = append(out, toView(p))
	}
	return &ListOutput{Products: out}, nil
}
