package application

import (
	"context"

	"github.com/google/uuid"
)

type GetInput struct {
	ID uuid.UUID
}

type GetOutput struct {
	Product ProductView
}

type GetProduct struct {
	repo Repository
}

func NewGetProduct(repo Repository) *GetProduct {
	return &GetProduct{repo: repo}
}

func (uc *GetProduct) Execute(ctx context.Context, in GetInput) (*GetOutput, error) {
	product, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	return &GetOutput{Product: toView(product)}, nil
}
