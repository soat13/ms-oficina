package application

import (
	"context"

	"github.com/google/uuid"
)

type DeleteVehicle struct {
	repo VehicleRepository
}

func NewDeleteVehicle(repo VehicleRepository) *DeleteVehicle {
	return &DeleteVehicle{repo: repo}
}

func (uc *DeleteVehicle) Execute(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
