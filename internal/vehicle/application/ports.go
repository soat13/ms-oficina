package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/vehicle/domain"
	"github.com/soat13/fase-1-oficina/pkg/utils/pagination"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/plate"
)

type (
	VehicleView struct {
		ID         uuid.UUID
		CustomerID uuid.UUID
		Plate      plate.Plate
		Brand      string
		Model      string
		Year       int
	}

	VehicleRepository interface {
		Create(ctx context.Context, v *domain.Vehicle) error
		Update(ctx context.Context, v *domain.Vehicle) error
		Delete(ctx context.Context, id uuid.UUID) error
		GetByID(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error)
		List(ctx context.Context, pager pagination.Pagination) ([]*domain.Vehicle, error)
		ListByCustomerID(ctx context.Context, customerID uuid.UUID, pager pagination.Pagination) ([]*domain.Vehicle, error)
		ExistsByPlate(ctx context.Context, plate plate.Plate) (bool, error)
	}
)

func toView(v *domain.Vehicle) VehicleView {
	if v == nil {
		return VehicleView{}
	}
	return VehicleView{
		ID:         v.ID,
		CustomerID: v.CustomerID,
		Plate:      v.Plate,
		Brand:      v.Brand,
		Model:      v.Model,
		Year:       v.Year,
	}
}
