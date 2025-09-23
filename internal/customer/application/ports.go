package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/customer/domain"
)

type CustomerView struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Cellphone    string    `json:"cellphone"`
	Document     string    `json:"document"`
	DocumentType string    `json:"document_type"`
}

func toView(c *domain.Customer) CustomerView {
	return CustomerView{
		ID:           c.ID,
		Name:         c.Name,
		Cellphone:    c.Cellphone,
		Document:     c.Document,
		DocumentType: c.DocumentType,
	}
}

type CustomerRepository interface {
	Create(ctx context.Context, s *domain.Customer) error
	Update(ctx context.Context, s *domain.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Customer, error)
	ExistsByDocument(ctx context.Context, document string) (bool, error)
}
