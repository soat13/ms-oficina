package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/oficina-utils/pkg/valueobjects/plate"
)

type (
	UpdateInput struct {
		ID    uuid.UUID
		Plate *plate.Plate
		Brand *string
		Model *string
		Year  *int
	}

	UpdateVehicle struct {
		repo VehicleRepository
	}
)

func NewUpdateVehicle(repo VehicleRepository) *UpdateVehicle {
	return &UpdateVehicle{repo: repo}
}

func (uv *UpdateVehicle) Execute(ctx context.Context, in UpdateInput) error {
	vehicle, err := uv.repo.GetByID(ctx, in.ID)
	if err != nil {
		return err
	}
	if vehicle == nil {
		return ErrVehicleNotFound
	}

	if in.Plate != nil {

		if in.Plate.String() != vehicle.Plate.String() {
			exists, err := uv.repo.ExistsByPlate(ctx, *in.Plate)
			if err != nil {
				return err
			}
			if exists {
				return ErrDuplicatePlate
			}
		}
		if err := vehicle.ChangePlate(*in.Plate); err != nil {
			return err
		}
	}
	if in.Brand != nil {
		if err := vehicle.ChangeBrand(*in.Brand); err != nil {
			return err
		}
	}
	if in.Model != nil {
		if err := vehicle.ChangeModel(*in.Model); err != nil {
			return err
		}
	}
	if in.Year != nil {
		if err := vehicle.ChangeYear(*in.Year); err != nil {
			return err
		}
	}

	return uv.repo.Update(ctx, vehicle)
}
