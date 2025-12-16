package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	GetInput struct {
		ID uuid.UUID
	}

	GetOutput struct {
		Service ServiceView
	}

	GetService struct {
		repo Repository
	}
)

func NewGetService(repo Repository) *GetService {
	return &GetService{repo: repo}
}

func (uc *GetService) Execute(ctx context.Context, in GetInput) (*GetOutput, error) {
	service, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, ErrServiceNotFound
	}
	return &GetOutput{Service: toView(service)}, nil
}
