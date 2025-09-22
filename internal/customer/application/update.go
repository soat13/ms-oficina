package application

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UpdateInput struct {
	ID        uuid.UUID
	Name      *string
	Cellphone *string
	Now       time.Time
}

type UpdateOutput struct {
	Customer CustomerView `json:"customer"`
}

type UpdateCustomer struct {
	repo CustomerRepository
}

func NewUpdateCustomer(repo CustomerRepository) *UpdateCustomer {
	return &UpdateCustomer{repo: repo}
}

func (uc *UpdateCustomer) Execute(ctx context.Context, in UpdateInput) (*UpdateOutput, error) {
	c, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCustomerNotFound
	}

	if in.Name != nil {
		if err := c.ChangeName(*in.Name, in.Now); err != nil {
			return nil, err
		}
	}
	if in.Cellphone != nil {
		if err := c.ChangeCellphone(*in.Cellphone, in.Now); err != nil {
			return nil, err
		}
	}

	if err := uc.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return &UpdateOutput{Customer: toView(c)}, nil
}
