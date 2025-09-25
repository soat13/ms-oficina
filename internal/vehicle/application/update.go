package application

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UpdateVehicle struct {
	repo VehicleRepository
}

func NewUpdateVehicle(repo VehicleRepository) *UpdateVehicle {
	return &UpdateVehicle{repo: repo}
}

type UpdateInput struct {
	ID    uuid.UUID
	Model string
	Brand string
	Year  int
}

func (uc *UpdateVehicle) Execute(ctx context.Context, in UpdateInput) error {
	vehicle, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return err
	}

	vehicle.Model = in.Model
	vehicle.Brand = in.Brand
	vehicle.Year = in.Year
	vehicle.UpdatedAt = time.Now()

	if err := vehicle.Validate(); err != nil {
		return err
	}

	return uc.repo.Update(ctx, vehicle)
}
