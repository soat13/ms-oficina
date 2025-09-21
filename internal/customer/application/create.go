package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/internal/customer/domain"
)

type CreateInput struct {
	Name      string
	Cellphone string
	Document  string
	Now       time.Time
}

type CustomerView struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Cellphone string    `json:"cellphone"`
	Document  string    `json:"document"`
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
	exists, err := cr.repo.ExistsByDocument(ctx, in.Document)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateCustomer
	}

	customer, err := domain.NewCustomer(uuid.Nil, in.Name, in.Cellphone, in.Document, in.Now)
	if err != nil {
		return nil, err
	}

	if err := cr.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return &CreateOutput{Customer: toView(customer)}, nil
}

func toView(c *domain.Customer) CustomerView {
	return CustomerView{
		ID:        c.ID,
		Name:      c.Name,
		Cellphone: c.Cellphone,
		Document:  c.Document,
	}
}
