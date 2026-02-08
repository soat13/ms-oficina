package application

import (
	"context"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
)

type FinishExecutionInput struct {
	RepairOrderID uuid.UUID
}

type FinishExecution struct {
	Repository       Repository
	metricsPublisher MetricsPublisher
}

func NewFinishExecution(repository Repository, metricsPublisher MetricsPublisher) *FinishExecution {
	return &FinishExecution{
		Repository:       repository,
		metricsPublisher: metricsPublisher,
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

	if err := uc.Repository.SaveIfInExecution(ctx, repairorder); err != nil {
		return err
	}

	uc.metricsPublisher.IncRepairOrderStatusChange("in_execution", "finished")

	return nil
}
