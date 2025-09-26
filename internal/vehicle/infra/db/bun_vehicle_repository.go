package db

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/uuid"
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

func (r *vehicleRepository) Update(ctx context.Context, vehicle *domain.Vehicle) error {
	_, err := r.db.NewUpdate().Model(vehicle).Where("id = ?", vehicle.ID).Exec(ctx)
	return err
}

func (r *vehicleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.Vehicle)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *vehicleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error) {
	vehicle := new(domain.Vehicle)
	err := r.db.NewSelect().Model(vehicle).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (r *vehicleRepository) List(ctx context.Context, limit, offset int) ([]*domain.Vehicle, error) {
	var vehicles []*domain.Vehicle
	err := r.db.NewSelect().Model(&vehicles).Limit(limit).Offset(offset).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (r vehicleRepository) ExistsByPlate(ctx context.Context, plate string) (bool, error) {
	plate = strings.TrimSpace(plate)

	// Padrão antigo: aaa-1234 ou aaa1234
	padraoAntigo := regexp.MustCompile(`^[a-z]{3}-?[0-9]{4}$`)

	// Padrão Mercosul: aaa1a23
	padraoMercosul := regexp.MustCompile(`^[a-z]{3}[0-9]{1}[a-z]{1}[0-9]{2}$`)

	valid_plate := padraoAntigo.MatchString(plate) || padraoMercosul.MatchString(plate)

	if valid_plate {
		exists, err := r.db.NewSelect().Model((*domain.Vehicle)(nil)).Where("plate = ?", plate).Exists(ctx)
		if err != nil {
			return false, err
		}
		return exists, nil
	}
	return false, nil
}
