package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/entity"
	uuidPkg "github.com/soat13/fase-1-oficina/pkg/utils/uuid"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/plate"
)

type Vehicle struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Plate      plate.Plate
	Brand      string
	Model      string
	Year       int
	entity.Timestamps
}

func NewVehicle(id uuid.UUID, customerID uuid.UUID, plateVO plate.Plate, brand, model string, year int) (*Vehicle, error) {
	brand = strings.TrimSpace(brand)
	if brand == "" {
		return nil, ErrInvalidVehicleBrand
	}

	model = strings.TrimSpace(model)
	if model == "" {
		return nil, ErrInvalidVehicleModel
	}

	if year < 1900 || year > time.Now().Year()+1 {
		return nil, ErrInvalidVehicleYear
	}

	now := time.Now()
	vehicle := &Vehicle{
		ID:         uuidPkg.IDOrNew(id),
		CustomerID: customerID,
		Plate:      plateVO,
		Brand:      brand,
		Model:      model,
		Year:       year,
		Timestamps: entity.NewTimestamps(now, now),
	}

	return vehicle, nil
}

func (v *Vehicle) ChangePlate(plate plate.Plate) error {
	v.Plate = plate
	v.Touch()
	return nil
}

func (v *Vehicle) ChangeBrand(newBrand string) error {
	newBrand = strings.TrimSpace(newBrand)
	if newBrand == "" {
		return ErrInvalidVehicleBrand
	}
	v.Brand = newBrand
	v.Touch()
	return nil
}

func (v *Vehicle) ChangeModel(newModel string) error {
	newModel = strings.TrimSpace(newModel)
	if newModel == "" {
		return ErrInvalidVehicleModel
	}
	v.Model = newModel
	v.Touch()
	return nil
}

func (v *Vehicle) ChangeYear(newYear int) error {
	if newYear < 1900 || newYear > time.Now().Year()+1 {
		return ErrInvalidVehicleYear
	}
	v.Year = newYear
	v.Touch()
	return nil
}
