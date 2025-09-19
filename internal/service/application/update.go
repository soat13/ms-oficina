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
	Currency   *string
	Now        time.Time
}

type UpdateOutput struct {
	Service ServiceView `json:"service"`
}

type UpdateService struct {
	repo Repository
}

func NewUpdateService(repo Repository) *UpdateService {
	return &UpdateService{repo: repo}
}

func (uc *UpdateService) Execute(ctx context.Context, in UpdateInput) (*UpdateOutput, error) {
	s, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, ErrServiceNotFound
	}

	if in.Name != nil {
		if err := s.Rename(*in.Name, in.Now); err != nil {
			return nil, err
		}
	}
	if in.PriceCents != nil {
		money, err := money.New(*in.PriceCents)
		if err != nil {
			return nil, err
		}

		if err := s.ChangePrice(money, in.Now); err != nil {
			return nil, err
		}
	}
	if in.Currency != nil {
		if err := s.ChangeCurrency(*in.Currency, in.Now); err != nil {
			return nil, err
		}
	}

	if err := uc.repo.Update(ctx, s); err != nil {
		return nil, err
	}
	return &UpdateOutput{Service: toView(s)}, nil
}
