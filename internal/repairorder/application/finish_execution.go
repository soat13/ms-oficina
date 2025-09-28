package application

import (
	"context"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
)

type FinishExecutionInput struct {
	RepairOrderID uuid.UUID
}

type FinishExecution struct {
	Repository Repository
}

func NewFinishExecution(repository Repository) *FinishExecution {
	return &FinishExecution{
		Repository: repository,
	}
}

func (uc *FinishExecution) Execute(ctx context.Context, input FinishExecutionInput) error {
	repairorder, err := uc.Repository.GetById(ctx, input.RepairOrderID)
	if err != nil {
		return err
	}

	if repairorder == nil {
		return sharedRepairOrder.ErrRepairOrderNotFound
	}

	if err := repairorder.FinishExecution(); err != nil {
		return err
	}

	return uc.Repository.SaveIfInExecution(ctx, repairorder)
}
