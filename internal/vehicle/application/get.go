package application

import (
	"context"

	"github.com/google/uuid"
)

type GetVehicle struct {
	repo VehicleRepository
}

func NewGetVehicle(repo VehicleRepository) *GetVehicle {
	return &GetVehicle{repo: repo}
}

func (uc *GetVehicle) Execute(ctx context.Context, id uuid.UUID) (*VehicleView, error) {
	vehicle, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	view := toView(vehicle)
	return &view, nil
}
