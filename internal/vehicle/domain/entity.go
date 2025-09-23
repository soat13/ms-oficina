package domain

import (
	"time"

	"github.com/google/uuid"
)

type Vehicle struct {
	ID         uuid.UUID `json:"id"`
	CustomerId string    `json:"customer_id"`
	Plate      string    `json:"plate"`
	Model      string    `json:"model"`
	Brand      string    `json:"brand"`
	Year       int       `json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewVehicle(customerId, plate, model, brand string, year int) (*Vehicle, error) {
	v := &Vehicle{
		ID:         uuid.New(),
		CustomerId: customerId,
		Plate:      plate,
		Model:      model,
		Brand:      brand,
		Year:       year,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := v.Validate(); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Vehicle) Validate() error {
	if v.CustomerId == "" {
		return ErrInvalidCustomerId
	}

	if v.Plate == "" {
		return ErrInvalidPlate
	}

	if v.Model == "" {
		return ErrInvalidModel
	}

	if v.Brand == "" {
		return ErrInvalidBrand
	}

	if v.Year <= 0 {
		return ErrInvalidYear
	}

	return nil
}

func (v *Vehicle) GetPlate() string {
	return v.Plate
}

func (v *Vehicle) GetBrand() string {
	return v.Brand
}

func (v *Vehicle) GetCustomerId() string {
	return v.CustomerId
}
