package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/product/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type CreateInput struct {
	Name  string
	Price money.Money
	Stock int // opcional: default 0
	Now   time.Time
}

type ProductView struct {
	ID         uuid.UUID   `json:"id"`
	Name       string      `json:"name"`
	Price      money.Money `json:"-"`
	PriceCents int64       `json:"price_cents"`
	Stock      int         `json:"stock"`
}

type CreateOutput struct {
	Product ProductView `json:"product"`
}

type CreateProduct struct {
	repo Repository
}

func NewCreateProduct(repo Repository) *CreateProduct {
	return &CreateProduct{repo: repo}
}

func (uc *CreateProduct) Execute(ctx context.Context, in CreateInput) (*CreateOutput, error) {
	exists, err := uc.repo.ExistsByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateProduct
	}

	p, err := domain.NewProduct(uuid.Nil, in.Name, in.Price, in.Stock, in.Now)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return &CreateOutput{Product: toView(p)}, nil
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
