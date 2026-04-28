package application

import (
	"context"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
)

type StartExecutionInput struct {
	RepairOrderID uuid.UUID
}

type StartExecution struct {
	Repository       Repository
	metricsPublisher MetricsPublisher
}

func NewStartExecution(repository Repository, metricsPublisher MetricsPublisher) *StartExecution {
	return &StartExecution{
		Repository:       repository,
		metricsPublisher: metricsPublisher,
	}
}

func (uc *StartExecution) Execute(ctx context.Context, input StartExecutionInput) error {
	repairorder, err := uc.Repository.GetById(ctx, input.RepairOrderID)
	if err != nil {
		return err
	}

	if repairorder == nil {
		return sharedRepairOrder.ErrRepairOrderNotFound
	}

	if err := repairorder.StartExecution(); err != nil {
		return err
	}

	return uc.Repository.SaveIfApproved(ctx, repairorder)
}
