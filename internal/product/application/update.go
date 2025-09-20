package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type UpdateInput struct {
	ID         uuid.UUID
	Name       *string
	PriceCents *int64
	Stock      *int
	Now        time.Time
}

type UpdateOutput struct {
	Product ProductView `json:"product"`
}

type UpdateProduct struct {
	repo Repository
}

func NewUpdateProduct(repo Repository) *UpdateProduct {
	return &UpdateProduct{repo: repo}
}

func (uc *UpdateProduct) Execute(ctx context.Context, in UpdateInput) (*UpdateOutput, error) {
	p, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProductNotFound
	}

	if in.Name != nil {
		if err := p.Rename(*in.Name, in.Now); err != nil {
			return nil, err
		}
	}
	if in.PriceCents != nil {
		m, err := money.New(*in.PriceCents)
		if err != nil {
			return nil, err
		}
		if err := p.ChangePrice(m, in.Now); err != nil {
			return nil, err
		}
	}
	if in.Stock != nil {
		if err := p.ChangeStock(*in.Stock, in.Now); err != nil {
			return nil, err
		}
	}

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return &UpdateOutput{Product: toView(p)}, nil
}
