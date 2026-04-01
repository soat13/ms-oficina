package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/oficina-utils/pkg/money"
)

type (
	UpdateInput struct {
		ID    uuid.UUID
		Name  *string
		Price *money.Money
	}

	UpdateService struct {
		repo Repository
	}
)

func NewUpdateService(repo Repository) *UpdateService {
	return &UpdateService{repo: repo}
}

func (uc *UpdateService) Execute(ctx context.Context, in UpdateInput) error {
	service, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return err
	}
	if service == nil {
		return ErrServiceNotFound
	}

	if in.Name != nil {
		if err := service.Rename(*in.Name); err != nil {
			return err
		}
	}

	if in.Price != nil {
		err := service.ChangePrice(*in.Price)
		if err != nil {
			return err
		}
	}

	return uc.repo.Update(ctx, service)
}
