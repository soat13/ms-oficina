package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/service/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type CreateInput struct {
	Name     string
	Price    money.Money
	Currency string // opcional: default BRL
	Now      time.Time
}

type ServiceView struct {
	ID         uuid.UUID   `json:"id"`
	Name       string      `json:"name"`
	Price      money.Money `json:"-"`
	PriceCents int64       `json:"price_cents"`
	Currency   string      `json:"currency"`
}

type CreateOutput struct {
	Service ServiceView `json:"service"`
}

type CreateService struct {
	repo Repository
}

func NewCreateService(repo Repository) *CreateService {
	return &CreateService{repo: repo}
}

func (uc *CreateService) Execute(ctx context.Context, in CreateInput) (*CreateOutput, error) {
	exists, err := uc.repo.ExistsByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateService
	}

	s, err := domain.NewService(uuid.Nil, in.Name, in.Price, in.Currency, in.Now)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Create(ctx, s); err != nil {
		return nil, err
	}
	return &CreateOutput{Service: toView(s)}, nil
}

func toView(s *domain.Service) ServiceView {
	return ServiceView{
		ID:         s.ID,
		Name:       s.Name,
		Price:      s.Price,
		PriceCents: s.Price.Cents,
		Currency:   s.Currency,
	}
}
