package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/internal/customer/domain"
	string_helper "github.com/soat13/fase-1-oficina/pkg/utils/helpers/string"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/document"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/phone"
)

type (
	CreateInput struct {
		Name        string
		Document    document.Document
		Email       email.Email
		PhoneNumber phone.PhoneNumber
	}

	CreateCustomer struct {
		repo CustomerRepository
	}
)

func NewCreateCustomer(repo CustomerRepository) *CreateCustomer {
	return &CreateCustomer{repo: repo}
}

func (cr *CreateCustomer) Execute(ctx context.Context, in CreateInput) error {
	err := checkIfCustomerExists(ctx, cr.repo, in.Document, in.Email)
	if err != nil {
		return err
	}

	customer, err := domain.NewCustomer(uuid.Nil, in.Name, in.Document, in.PhoneNumber, in.Email)
	if err != nil {
		return err
	}

	return cr.repo.Create(ctx, customer)
}

func checkIfCustomerExists(ctx context.Context, repo CustomerRepository, document document.Document, email email.Email) error {
	existsDoc, err := repo.ExistsByDocument(ctx, string_helper.OnlyNumbers(document.Value))
	if err != nil {
		return err
	}
	if existsDoc {
		return ErrDuplicateDocument
	}

	existsEmail, err := repo.ExistsByEmail(ctx, email.String())
	if err != nil {
		return err
	}
	if existsEmail {
		return ErrDuplicateEmail
	}

	return nil
}
