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
	Stock int
	Now   time.Time
}

type CreateOutput struct {
	Product ProductView
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

	product, err := domain.NewProduct(uuid.Nil, in.Name, in.Price, in.Stock, in.Now, in.Now)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return &CreateOutput{Product: toView(product)}, nil
}
