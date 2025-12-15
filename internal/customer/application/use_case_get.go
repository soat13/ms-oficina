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
		Customer CustomerView
	}

	GetCustomer struct {
		repo CustomerRepository
	}
)

func NewGetCustomer(repo CustomerRepository) *GetCustomer {
	return &GetCustomer{repo: repo}
}

func (uc *GetCustomer) Execute(ctx context.Context, in GetInput) (*GetOutput, error) {
	customer, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}
	return &GetOutput{Customer: toView(customer)}, nil
}
