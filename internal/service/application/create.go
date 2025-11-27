package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/service/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type (
	CreateInput struct {
		Name  string
		Price money.Money
		Now   time.Time
	}

	CreateOutput struct {
		ServiceID uuid.UUID
	}

	CreateService struct {
		repo ServiceRepository
	}
)

func NewCreateService(repo ServiceRepository) *CreateService {
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

	service, err := domain.NewService(uuid.Nil, input.Name, input.Price, input.Now, input.Now)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Create(ctx, service)
	if err != nil {
		return nil, err
	}

	return &CreateOutput{ServiceID: service.ID}, nil
}
