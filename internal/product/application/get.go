package application

import (
	"context"

	"github.com/google/uuid"
)

type GetInput struct {
	ID uuid.UUID
}

type GetOutput struct {
	Product ProductView `json:"product"`
}

type GetProduct struct {
	repo ProductRepository
}

func NewGetProduct(repo ProductRepository) *GetProduct {
	return &GetProduct{repo: repo}
}

func (uc *GetProduct) Execute(ctx context.Context, in GetInput) (*GetOutput, error) {
	p, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProductNotFound
	}
	return &GetOutput{Product: toView(p)}, nil
}
