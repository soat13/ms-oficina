package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/entity"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
)

type BunRepairOrderRepository struct {
	db *bun.DB
}

func NewBunRepairOrderRepository(db *bun.DB) application.RepairOrderRepository {
	return &BunRepairOrderRepository{db: db}
}

type repairOrderModel struct {
	bun.BaseModel `bun:"table:repair_orders"`

	ID         uuid.UUID          `bun:"id,pk"`
	CustomerID uuid.UUID          `bun:"customer_id"`
	VehicleID  uuid.UUID          `bun:"vehicle_id"`
	Status     repairorder.Status `bun:"status"`
	CreatedAt  time.Time          `bun:"created_at"`
	UpdatedAt  time.Time          `bun:"updated_at"`
}

func (r *BunRepairOrderRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.RepairOrder, error) {
	var m repairOrderModel
	err := r.db.NewSelect().
		Model(&m).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &domain.RepairOrder{
		ID:         m.ID,
		CustomerID: m.CustomerID,
		VehicleID:  m.VehicleID,
		Status:     m.Status,
		Timestamps: entity.Timestamps{
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
	}, nil
}

func (r *BunRepairOrderRepository) Save(ctx context.Context, repairOrder *domain.RepairOrder) (*domain.RepairOrder, error) {
	m := repairOrderModel{
		ID:         repairOrder.ID,
		CustomerID: repairOrder.CustomerID,
		VehicleID:  repairOrder.VehicleID,
		Status:     repairOrder.Status,
		CreatedAt:  repairOrder.CreatedAt,
		UpdatedAt:  repairOrder.UpdatedAt,
	}

	_, err := r.db.NewInsert().
		Model(&m).
		On("CONFLICT (id) DO UPDATE").
		Set("customer_id = EXCLUDED.customer_id").
		Set("vehicle_id = EXCLUDED.vehicle_id").
		Set("status = EXCLUDED.status").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return repairOrder, nil
}
