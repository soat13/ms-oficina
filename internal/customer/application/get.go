package application

import (
	"context"

	"github.com/google/uuid"
)

type GetInput struct {
	ID uuid.UUID
}

type GetOutput struct {
	Customer CustomerView `json:"customer"`
}

type GetCustomer struct {
	repo CustomerRepository
}

func NewGetCustomer(repo CustomerRepository) *GetCustomer {
	return &GetCustomer{repo: repo}
}

func (uc *GetCustomer) Execute(ctx context.Context, in GetInput) (*GetOutput, error) {
	c, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCustomerNotFound
	}
	return &GetOutput{Customer: toView(c)}, nil
}
