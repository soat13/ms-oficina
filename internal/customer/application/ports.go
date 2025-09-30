package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/customer/domain"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/document"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/phone"
)

type CustomerView struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Document    document.Document `json:"document"`
	Email       email.Email       `json:"email"`
	PhoneNumber phone.PhoneNumber `json:"phone_number"`
}

func toView(c *domain.Customer) CustomerView {
	return CustomerView{
		ID:          c.ID,
		Name:        c.Name,
		Document:    c.Document,
		Email:       c.Email,
		PhoneNumber: c.PhoneNumber,
	}
}

type CustomerRepository interface {
	Create(ctx context.Context, s *domain.Customer) error
	Update(ctx context.Context, s *domain.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Customer, error)
	ExistsByDocument(ctx context.Context, document string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
