package application

import (
	"context"

	"github.com/google/uuid"
)

type DeleteInput struct {
	ID uuid.UUID
}

type DeleteVehicle struct {
	repo VehicleRepository
}

func NewDeleteVehicle(repo VehicleRepository) *DeleteVehicle {
	return &DeleteVehicle{repo: repo}
}

func (uc *DeleteVehicle) Execute(ctx context.Context, in DeleteInput) error {
	return uc.repo.Delete(ctx, in.ID)
}
