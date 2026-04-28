package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
)

type FinishExecutionInput struct {
	RepairOrderID uuid.UUID
}

type FinishExecution struct {
	Repository       Repository
	eventPublisher   EventPublisher
	metricsPublisher MetricsPublisher
}

func NewFinishExecution(repository Repository, eventPublisher EventPublisher, metricsPublisher MetricsPublisher) *FinishExecution {
	return &FinishExecution{
		Repository:       repository,
		eventPublisher:   eventPublisher,
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

	if repairorder.TotalEstimate == nil {
		return ErrTotalEstimateRequired
	}

	if err := repairorder.FinishExecution(); err != nil {
		return err
	}

	if err := uc.Repository.SaveIfInExecution(ctx, repairorder); err != nil {
		return err
	}

	if repairorder.Timestamps != nil {
		uc.metricsPublisher.RecordRepairOrderPhaseDuration("in_execution", time.Since(repairorder.UpdatedAt).Minutes())
	}

	if err := uc.eventPublisher.PublishPaymentRequest(ctx, repairorder); err != nil {
		return err
	}

	if err := repairorder.RequestPayment(); err != nil {
		return err
	}

	return uc.Repository.SaveIfFinished(ctx, repairorder)
}
