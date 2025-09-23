package db

import (
	"context"

	"github.com/soat13/fase-1-oficina/internal/vehicle/application"
	"github.com/soat13/fase-1-oficina/internal/vehicle/domain"
	"github.com/uptrace/bun"
)

type vehicleRepository struct {
	db *bun.DB
}

func NewVehicleRepository(db *bun.DB) application.VehicleRepository {
	return &vehicleRepository{
		db: db,
	}
}

func (r *vehicleRepository) Create(ctx context.Context, vehicle *domain.Vehicle) (*domain.Vehicle, error) {
	_, err := r.db.NewInsert().Model(vehicle).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (r vehicleRepository) ExistsByPlate(ctx context.Context, plate string) (bool, error) {
	exists, err := r.db.NewSelect().Model((*domain.Vehicle)(nil)).Where("plate = ?", plate).Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}
