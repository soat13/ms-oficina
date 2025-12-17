package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/vehicle/application"
	"github.com/soat13/fase-1-oficina/internal/vehicle/domain"
	"github.com/soat13/fase-1-oficina/pkg/db/bun_helper"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/pagination"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/plate"
	"github.com/uptrace/bun"
)

type vehicleModel struct {
	bun.BaseModel `bun:"table:vehicles"`

	ID         uuid.UUID `bun:",pk,type:uuid"`
	CustomerID uuid.UUID `bun:",notnull,type:uuid"`
	Plate      string    `bun:",notnull"`
	Brand      string    `bun:",notnull"`
	Model      string    `bun:",notnull"`
	Year       int       `bun:",notnull"`
	CreatedAt  time.Time `bun:",nullzero,default:now()"`
	UpdatedAt  time.Time `bun:",nullzero,default:now()"`
}

type BunVehicleRepository struct {
	db *bun.DB
}

func NewBunVehicleRepository(db *bun.DB) application.VehicleRepository {
	return &BunVehicleRepository{db: db}
}

func (repo *BunVehicleRepository) Create(ctx context.Context, vehicle *domain.Vehicle) error {
	model := toModel(vehicle)
	_, err := repo.db.NewInsert().Model(model).Exec(ctx)
	return err
}

func (repo *BunVehicleRepository) Update(ctx context.Context, vehicle *domain.Vehicle) error {
	model := toModel(vehicle)
	_, err := repo.db.NewUpdate().
		Model(model).
		Column("plate", "brand", "model", "year", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (repo *BunVehicleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.NewDelete().Model(&vehicleModel{ID: id}).WherePK().Exec(ctx)
	return bun_helper.HandleDeleteError(err)
}

func (repo *BunVehicleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error) {
	model := vehicleModel{ID: id}
	err := repo.db.NewSelect().Model(&model).WherePK().Scan(ctx)

	if err := bun_helper.IgnoreNoRows(err); err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (repo *BunVehicleRepository) List(ctx context.Context, pager pagination.Pagination) ([]*domain.Vehicle, error) {
	var rows []vehicleModel
	if err := repo.db.NewSelect().
		Model(&rows).
		Order("brand ASC").
		Order("model ASC").
		Limit(pager.Limit).
		Offset(pager.Offset).
		Scan(ctx); err != nil {
		return nil, err
	}

	return maps.MapPtr(rows, toDomain), nil
}

func (repo *BunVehicleRepository) ListByCustomerID(ctx context.Context, customerID uuid.UUID, pager pagination.Pagination) ([]*domain.Vehicle, error) {
	var rows []vehicleModel
	if err := repo.db.NewSelect().
		Model(&rows).
		Where("customer_id = ?", customerID).
		Order("brand ASC").
		Order("model ASC").
		Limit(pager.Limit).
		Offset(pager.Offset).
		Scan(ctx); err != nil {
		return nil, err
	}

	return maps.MapPtr(rows, toDomain), nil
}

func (repo *BunVehicleRepository) ExistsByPlate(ctx context.Context, plateVO plate.Plate) (bool, error) {
	return repo.db.NewSelect().
		Model((*vehicleModel)(nil)).
		Where("plate = ?", plateVO.String()).
		Exists(ctx)
}

func toModel(v *domain.Vehicle) *vehicleModel {
	return &vehicleModel{
		ID:         v.ID,
		CustomerID: v.CustomerID,
		Plate:      v.Plate.String(),
		Brand:      v.Brand,
		Model:      v.Model,
		Year:       v.Year,
		CreatedAt:  v.CreatedAt,
		UpdatedAt:  v.UpdatedAt,
	}
}

func toDomain(model *vehicleModel) *domain.Vehicle {
	if model == nil || model.ID == uuid.Nil {
		return nil
	}

	plateVO, _ := plate.New(model.Plate)

	vehicle, err := domain.NewVehicle(model.ID, model.CustomerID, plateVO, model.Brand, model.Model, model.Year)
	if err != nil {
		return nil
	}

	vehicle.UpdatedAt = model.UpdatedAt
	return vehicle
}
