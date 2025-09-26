package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/product/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type ProductView struct {
	ID         uuid.UUID   `json:"id"`
	Name       string      `json:"name"`
	Price      money.Money `json:"-"`
	PriceCents int64       `json:"price_cents"`
	Stock      int         `json:"stock"`
}

func toView(p *domain.Product) ProductView {
	return ProductView{
		ID:         p.ID,
		Name:       p.Name,
		Price:      p.Price,
		PriceCents: p.Price.Cents,
		Stock:      p.Stock,
	}
}

type Repository interface {
	Create(ctx context.Context, p *domain.Product) error
	Update(ctx context.Context, p *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Product, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}
