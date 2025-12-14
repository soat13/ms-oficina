package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/product/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type (
	CreateInput struct {
		Name  string
		Price money.Money
		Stock int
	}

	CreateOutput struct {
		ProductID uuid.UUID
	}

	CreateProduct struct {
		repo Repository
	}
)

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

	product, err := domain.NewProduct(uuid.Nil, in.Name, in.Price, in.Stock)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return &CreateOutput{ProductID: product.ID}, nil
}
