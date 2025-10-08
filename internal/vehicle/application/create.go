package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/vehicle/domain"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/plate"
)

type (
	CreateInput struct {
		CustomerID uuid.UUID
		Plate      string
		Brand      string
		Model      string
		Year       int
	}

	CreateVehicle struct {
		repo VehicleRepository
	}
)

func NewCreateVehicle(repo VehicleRepository) *CreateVehicle {
	return &CreateVehicle{repo: repo}
}

func (cv *CreateVehicle) Execute(ctx context.Context, in CreateInput) error {
	plateVO, err := plate.New(in.Plate)
	if err != nil {
		return err
	}

	exists, err := cv.repo.ExistsByPlate(ctx, plateVO)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicatePlate
	}

	vehicle, err := domain.NewVehicle(uuid.Nil, in.CustomerID, plateVO, in.Brand, in.Model, in.Year)
	if err != nil {
		return err
	}

	return cv.repo.Create(ctx, vehicle)
}
