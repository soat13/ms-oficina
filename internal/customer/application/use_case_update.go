package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/oficina-utils/pkg/valueobjects/email"
	"github.com/soat13/oficina-utils/pkg/valueobjects/phone"
)

type (
	UpdateInput struct {
		ID          uuid.UUID
		Name        *string
		PhoneNumber *phone.PhoneNumber
		Email       *email.Email
	}

	UpdateCustomer struct {
		repo CustomerRepository
	}
)

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
		if err := customer.ChangeName(*in.Name); err != nil {
			return err
		}
	}
	if in.PhoneNumber != nil {
		customer.ChangePhoneNumber(*in.PhoneNumber)
	}
	if in.Email != nil && customer.Email.String() != in.Email.String() {
		exists, err := uc.repo.ExistsByEmail(ctx, in.Email.String())
		if err != nil {
			return err
		}
		if !exists {
			customer.ChangeEmail(*in.Email)
		}
	}

	return uc.repo.Update(ctx, customer)
}
