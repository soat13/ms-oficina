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
	limit, offset := limitAndOffset(in)
	items, err := uc.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]ServiceView, 0, len(items))
	for _, s := range items {
		out = append(out, toView(s))
	}
	return &ListOutput{Services: out}, nil
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
