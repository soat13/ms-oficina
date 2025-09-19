package application

import (
	"context"

	"github.com/google/uuid"
)

type GetInput struct {
	ID uuid.UUID
}

type GetOutput struct {
	Service ServiceView `json:"service"`
}

type GetService struct {
	repo Repository
}

func NewGetService(repo Repository) *GetService {
	return &GetService{repo: repo}
}

func (uc *GetService) Execute(ctx context.Context, in GetInput) (*GetOutput, error) {
	s, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, ErrServiceNotFound
	}
	return &GetOutput{Service: toView(s)}, nil
}
