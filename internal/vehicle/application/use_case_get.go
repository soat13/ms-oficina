package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	GetInput struct {
		ID uuid.UUID
	}

	GetOutput struct {
		Vehicle VehicleView
	}

	GetVehicle struct {
		repo VehicleRepository
	}
)

func NewGetVehicle(repo VehicleRepository) *GetVehicle {
	return &GetVehicle{repo: repo}
}

func (gv *GetVehicle) Execute(ctx context.Context, in GetInput) (GetOutput, error) {
	vehicle, err := gv.repo.GetByID(ctx, in.ID)
	if err != nil {
		return GetOutput{}, err
	}
	if vehicle == nil {
		return GetOutput{}, ErrVehicleNotFound
	}

	return GetOutput{
		Vehicle: toView(vehicle),
	}, nil
}
