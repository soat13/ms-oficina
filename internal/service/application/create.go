package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/service/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type CreateInput struct {
	Name  string
	Price money.Money
	Now   time.Time
}

type ServiceView struct {
	ID    uuid.UUID   `json:"id"`
	Name  string      `json:"name"`
	Price money.Money `json:"price"`
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

func (uc *CreateService) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	exists, err := uc.repo.ExistsByName(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateService
	}

	s, err := domain.NewService(uuid.Nil, input.Name, input.Price, input.Now, input.Now)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, s); err != nil {
		return nil, err
	}

	return &CreateOutput{Service: toView(s)}, nil
}
