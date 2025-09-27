package application

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UpdateInput struct {
	ID          uuid.UUID
	Name        *string
	PhoneNumber *string
	Email       *string
	Now         time.Time
}

type UpdateCustomer struct {
	repo CustomerRepository
}

func NewUpdateCustomer(repo CustomerRepository) *UpdateCustomer {
	return &UpdateCustomer{repo: repo}
}

func (uc *UpdateCustomer) Execute(ctx context.Context, in UpdateInput) error {
	customer, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return err
	}
	if customer == nil {
		return ErrCustomerNotFound
	}
	if in.Name != nil {
		if err := customer.ChangeName(*in.Name, in.Now); err != nil {
			return err
		}
	}
	if in.PhoneNumber != nil {
		if err := customer.ChangePhoneNumber(*in.PhoneNumber, in.Now); err != nil {
			return err
		}
	}
	if in.Email != nil {
		exists, err := uc.repo.ExistsByEmail(ctx, *in.Email)
		if err != nil {
			return err
		}
		if exists {
			return ErrDuplicateCustomer
		}

		if err := customer.ChangeEmail(*in.Email, in.Now); err != nil {
			return err
		}
	}

	return uc.repo.Update(ctx, customer)
}
