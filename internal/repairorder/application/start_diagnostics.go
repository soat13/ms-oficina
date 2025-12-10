package application

import (
	"context"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
)

type StartDiagnosticsInput struct {
	RepairOrderID uuid.UUID
}

type StartDiagnostics struct {
	Repository Repository
}

func NewStartDiagnostics(repository Repository) *StartDiagnostics {
	return &StartDiagnostics{
		Repository: repository,
	}
}

func (uc *StartDiagnostics) Execute(ctx context.Context, input StartDiagnosticsInput) error {
	repairorder, err := uc.Repository.GetById(ctx, input.RepairOrderID)
	if err != nil {
		return err
	}
	if repairorder == nil {
		return sharedRepairOrder.ErrRepairOrderNotFound
	}

	if err := repairorder.StartDiagnostics(); err != nil {
		return err
	}

	return uc.Repository.SaveIfReceived(ctx, repairorder)
}
