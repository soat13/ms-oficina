package application

import (
	"context"

	"github.com/google/uuid"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
)

type ReleaseVehicleInput struct {
	RepairOrderID uuid.UUID
}

type ReleaseVehicle struct {
	Repository Repository
}

func NewReleaseVehicle(repository Repository) *ReleaseVehicle {
	return &ReleaseVehicle{
		Repository: repository,
	}
}

func (uc *ReleaseVehicle) Execute(ctx context.Context, input ReleaseVehicleInput) error {
	repairorder, err := uc.Repository.GetById(ctx, input.RepairOrderID)
	if err != nil {
		return err
	}

	if repairorder == nil {
		return sharedRepairOrder.ErrRepairOrderNotFound
	}

	if err := repairorder.ReleaseVehicle(); err != nil {
		return err
	}

	return uc.Repository.SaveIfFinished(ctx, repairorder)
}
