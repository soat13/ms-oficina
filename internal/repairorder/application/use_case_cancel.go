package application

import (
	"context"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
)

type CancelInput struct {
	RepairOrderID uuid.UUID
}

type Cancel struct {
	Repository       Repository
	eventPublisher   EventPublisher
	metricsPublisher MetricsPublisher
}

func NewCancel(repository Repository, eventPublisher EventPublisher, metricsPublisher MetricsPublisher) *Cancel {
	return &Cancel{
		Repository:       repository,
		eventPublisher:   eventPublisher,
		metricsPublisher: metricsPublisher,
	}
}

func (uc *Cancel) Execute(ctx context.Context, input CancelInput) error {
	repairorder, err := uc.Repository.GetById(ctx, input.RepairOrderID)
	if err != nil {
		return err
	}

	if repairorder == nil {
		return sharedRepairOrder.ErrRepairOrderNotFound
	}

	if repairorder.IsCancelled() {
		return nil
	}

	if err := repairorder.Cancel(); err != nil {
		return err
	}

	if err := uc.Repository.SaveCancellation(ctx, repairorder); err != nil {
		return err
	}

	uc.metricsPublisher.IncRepairOrderCanceled()

	return uc.eventPublisher.PublishRepairOrderCanceled(ctx, repairorder.ID)
}
