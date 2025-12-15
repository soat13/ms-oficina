package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	DeleteInput struct {
		ID uuid.UUID
	}

	DeleteCustomer struct {
		repo CustomerRepository
	}
)

func NewDeleteCustomer(repo CustomerRepository) *DeleteCustomer {
	return &DeleteCustomer{repo: repo}
}

func (uc *DeleteCustomer) Execute(ctx context.Context, in DeleteInput) error {
	return uc.repo.Delete(ctx, in.ID)
}
