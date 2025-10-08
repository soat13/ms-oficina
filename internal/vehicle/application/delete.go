package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	DeleteInput struct {
		ID uuid.UUID
	}

	DeleteVehicle struct {
		repo VehicleRepository
	}
)

func NewDeleteVehicle(repo VehicleRepository) *DeleteVehicle {
	return &DeleteVehicle{repo: repo}
}

func (dv *DeleteVehicle) Execute(ctx context.Context, in DeleteInput) error {
	vehicle, err := dv.repo.GetByID(ctx, in.ID)
	if err != nil {
		return err
	}
	if vehicle == nil {
		return ErrVehicleNotFound
	}

	return dv.repo.Delete(ctx, in.ID)
}
