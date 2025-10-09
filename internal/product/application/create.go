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

	CreateProduct struct {
		repo ProductRepository
	}
)

func NewCreateProduct(repo ProductRepository) *CreateProduct {
	return &CreateProduct{repo: repo}
}

func (uc *CreateProduct) Execute(ctx context.Context, in CreateInput) error {
	exists, err := uc.repo.ExistsByName(ctx, in.Name)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicateProduct
	}

	product, err := domain.NewProduct(uuid.Nil, in.Name, in.Price, in.Stock)
	if err != nil {
		return err
	}

	return uc.repo.Create(ctx, product)
}
