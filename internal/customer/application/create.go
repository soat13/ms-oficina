package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/internal/customer/domain"
	string_helper "github.com/soat13/fase-1-oficina/pkg/utils/helpers/string"
)

type CreateInput struct {
	Name        string
	Document    string
	Email       string
	PhoneNumber string
	Now         time.Time
}

type CreateCustomer struct {
	repo CustomerRepository
}

func NewCreateCustomer(repo CustomerRepository) *CreateCustomer {
	return &CreateCustomer{repo: repo}
}

func (cr *CreateCustomer) Execute(ctx context.Context, in CreateInput) error {
	exists, err := checkIfCustomerExists(ctx, cr.repo, in.Document, in.Email)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicateCustomer
	}

	customer, err := domain.NewCustomer(uuid.Nil, in.Name, in.Document, in.PhoneNumber, in.Email, in.Now)
	if err != nil {
		return err
	}

	return cr.repo.Create(ctx, customer)
}

func checkIfCustomerExists(ctx context.Context, repo CustomerRepository, document, email string) (bool, error) {
	document = string_helper.OnlyNumbers(document)
	existsDoc, err := repo.ExistsByDocument(ctx, document)
	if err != nil {
		return false, err
	}

	existsEmail, err := repo.ExistsByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return existsDoc || existsEmail, nil
}
