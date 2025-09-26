package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/internal/customer/domain"
)

type CreateInput struct {
	Name      string
	Document  string
	Cellphone string
	Email     string
	Now       time.Time
}

type CreateOutput struct {
	Customer CustomerView `json:"customer"`
}

type CreateCustomer struct {
	repo CustomerRepository
}

func NewCreateCustomer(repo CustomerRepository) *CreateCustomer {
	return &CreateCustomer{repo: repo}
}

func (cr *CreateCustomer) Execute(ctx context.Context, in CreateInput) (*CreateOutput, error) {
	exists, err := checkIfCustomerExists(ctx, cr.repo, in.Document, in.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateCustomer
	}

	customer, err := domain.NewCustomer(uuid.Nil, in.Name, in.Document, in.Cellphone, in.Email, in.Now)
	if err != nil {
		return nil, err
	}

	if err := cr.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return &CreateOutput{Customer: toView(customer)}, nil
}

func checkIfCustomerExists(ctx context.Context, repo CustomerRepository, document, email string) (bool, error) {
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
