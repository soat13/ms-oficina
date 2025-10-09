package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/pkg/money"
)

type (
	UpdateInput struct {
		ID    uuid.UUID
		Name  *string
		Price *money.Money
		Stock *int
	}

	UpdateProduct struct {
		repo ProductRepository
	}
)

func NewUpdateProduct(repo ProductRepository) *UpdateProduct {
	return &UpdateProduct{repo: repo}
}

func (uc *UpdateProduct) Execute(ctx context.Context, in UpdateInput) error {
	product, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}

	if in.Name != nil {
		if err := product.Rename(*in.Name); err != nil {
			return err
		}
	}

	if in.Price != nil {
		if err := product.ChangePrice(*in.Price); err != nil {
			return err
		}
	}

	if in.Stock != nil {
		if err := product.UpdateStock(*in.Stock); err != nil {
			return err
		}
	}

	return uc.repo.Update(ctx, product)
}
