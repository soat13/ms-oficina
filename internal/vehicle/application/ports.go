package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/vehicle/domain"
)

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
