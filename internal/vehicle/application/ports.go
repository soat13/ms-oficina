package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/vehicle/domain"
)

type VehicleView struct {
	ID         uuid.UUID `json:"id"`
	CustomerId string    `json:"customer_id"`
	Plate      string    `json:"plate"`
	Model      string    `json:"model"`
	Brand      string    `json:"brand"`
	Year       int       `json:"year"`
}

func toView(s *domain.Vehicle) VehicleView {
	return VehicleView{
		ID:         s.ID,
		CustomerId: s.CustomerId,
		Plate:      s.Plate,
		Model:      s.Model,
		Brand:      s.Brand,
		Year:       s.Year,
	}
}

type (
	VehicleRepository interface {
		Create(ctx context.Context, vehicle *domain.Vehicle) (*domain.Vehicle, error)
		Update(ctx context.Context, s *domain.Vehicle) error
		Delete(ctx context.Context, id uuid.UUID) error
		GetByID(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error)
		List(ctx context.Context, limit, offset int) ([]*domain.Vehicle, error)
		ExistsByPlate(ctx context.Context, plate string) (bool, error)
	}
)
