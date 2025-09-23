package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/vehicle/domain"
)

type CreateInput struct {
	CustomerId string
	Plate      string
	Model      string
	Brand      string
	Year       int
}

type VehicleView struct {
	ID         uuid.UUID `json:"id"`
	CustomerId string    `json:"customer_id"`
	Plate      string    `json:"plate"`
	Model      string    `json:"model"`
	Brand      string    `json:"brand"`
	Year       int       `json:"year"`
}

type CreateOutput struct {
	Vehicle VehicleView `json:"vehicle"`
}

type CreateVehicle struct {
	repo VehicleRepository
}

func NewCreateVehicle(repo VehicleRepository) *CreateVehicle {
	return &CreateVehicle{repo: repo}
}

func (uc *CreateVehicle) Execute(ctx context.Context, in CreateInput) (*CreateOutput, error) {
	exists, err := uc.repo.ExistsByPlate(ctx, in.Plate)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateVehicle
	}

	s, err := domain.NewVehicle(in.CustomerId, in.Plate, in.Model, in.Brand, in.Year)
	if err != nil {
		return nil, err
	}
	if _, err := uc.repo.Create(ctx, s); err != nil {
		return nil, err
	}
	return &CreateOutput{Vehicle: toView(s)}, nil
}

func toView(s *domain.Vehicle) VehicleView {
	return VehicleView{
		CustomerId: s.CustomerId,
		Plate:      s.Plate,
		Model:      s.Model,
		Brand:      s.Brand,
		Year:       s.Year,
	}
}
