package application

import (
	"context"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
)

type StartExecutionInput struct {
	RepairOrderID uuid.UUID
}

type StartExecution struct {
	Repository Repository
}

type RepairOrderExecutionStarted struct {
	OrderID    string
	ActorID    string
	OccurredAt int64
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
