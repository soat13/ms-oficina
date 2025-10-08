package application

import (
	"context"

	"github.com/google/uuid"
)

type DeleteInput struct {
	ID uuid.UUID
}

type DeleteProduct struct {
	repo Repository
}

func NewDeleteProduct(repo Repository) *DeleteProduct {
	return &DeleteProduct{repo: repo}
}

func (uc *DeleteProduct) Execute(ctx context.Context, in DeleteInput) error {
	return uc.repo.Delete(ctx, in.ID)
}
