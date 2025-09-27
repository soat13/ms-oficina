package application

import (
	"context"

	"github.com/google/uuid"
)

type GetVehicle struct {
	repo VehicleRepository
}

type GetInput struct {
	ID uuid.UUID
}

type GetOutput struct {
	Vehicle VehicleView `json:"vehicle"`
}

func NewGetVehicle(repo VehicleRepository) *GetVehicle {
	return &GetVehicle{repo: repo}
}

func (uc *GetVehicle) Execute(ctx context.Context, in GetInput) (*GetOutput, error) {
	vehicle, err := uc.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if vehicle == nil {
		return nil, ErrVehicleNotFound
	}

	return &GetOutput{Vehicle: toView(vehicle)}, nil
}
