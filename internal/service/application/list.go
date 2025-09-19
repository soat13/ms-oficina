package application

import "context"

type ListInput struct {
	Limit  int
	Offset int
}

type ListOutput struct {
	Services []ServiceView `json:"services"`
}

type ListServices struct {
	repo Repository
}

func NewListServices(repo Repository) *ListServices {
	return &ListServices{repo: repo}
}

func (uc *ListServices) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
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
	out := make([]ServiceView, 0, len(items))
	for _, s := range items {
		out = append(out, toView(s))
	}
	return &ListOutput{Services: out}, nil
}
