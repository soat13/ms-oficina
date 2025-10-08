package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/pkg/money"
)

type UpdateInput struct {
	ID    uuid.UUID
	Name  *string
	Price *money.Money
	Stock *int
	Now   time.Time
}

type UpdateOutput struct {
	Product ProductView
}

type UpdateProduct struct {
	repo Repository
}

func NewUpdateProduct(repo Repository) *UpdateProduct {
	return &UpdateProduct{repo: repo}
}

func (uc *UpdateProduct) Execute(ctx context.Context, in UpdateInput) (*UpdateOutput, error) {
	product, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	if in.Name != nil {
		if err := product.Rename(*in.Name, in.Now); err != nil {
			return nil, err
		}
	}

	if in.Price != nil {
		if err := product.ChangePrice(*in.Price); err != nil {
			return nil, err
		}
	}

	if in.Stock != nil {
		if err := product.UpdateStock(*in.Stock, in.Now); err != nil {
			return nil, err
		}
	}

	if err := uc.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return &UpdateOutput{Product: toView(product)}, nil
}
