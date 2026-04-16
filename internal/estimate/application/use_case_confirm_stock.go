package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	ConfirmStockInput struct {
		EstimateID uuid.UUID
	}

	ConfirmStock struct {
		repository Repository
	}
)

func NewConfirmStock(repository Repository) *ConfirmStock {
	return &ConfirmStock{repository: repository}
}

func (uc *ConfirmStock) Execute(ctx context.Context, input ConfirmStockInput) error {
	estimate, err := uc.repository.GetByID(ctx, input.EstimateID)
	if err != nil {
		return err
	}
	if estimate == nil {
		return ErrEstimateNotFound
	}

	if err := estimate.ConfirmStock(); err != nil {
		return err
	}

	return uc.repository.SaveIfAwaitingStock(ctx, estimate)
}
