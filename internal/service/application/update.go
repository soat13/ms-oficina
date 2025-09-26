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
	Now   time.Time
}

type UpdateOutput struct {
	Service ServiceView `json:"service"`
}

type UpdateService struct {
	repo ServiceRepository
}

func NewUpdateService(repo ServiceRepository) *UpdateService {
	return &UpdateService{repo: repo}
}

func (uc *UpdateService) Execute(ctx context.Context, in UpdateInput) (*UpdateOutput, error) {

	service, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, ErrServiceNotFound
	}

	if in.Name != nil {
		if err := service.Rename(*in.Name, in.Now); err != nil {
			return nil, err
		}
	}

	if in.Price != nil {
		if err := service.ChangePrice(*in.Price); err != nil {
			return nil, err
		}
	}

	if err := uc.repo.Update(ctx, service); err != nil {
		return nil, err
	}

	return &UpdateOutput{Service: toView(service)}, nil
}
