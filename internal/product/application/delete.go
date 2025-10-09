package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	DeleteInput struct {
		ID uuid.UUID
	}

	DeleteProduct struct {
		repo ProductRepository
	}
)

func NewDeleteProduct(repo ProductRepository) *DeleteProduct {
	return &DeleteProduct{repo: repo}
}

func (uc *DeleteProduct) Execute(ctx context.Context, in DeleteInput) error {
	return uc.repo.Delete(ctx, in.ID)
}
