package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/shared/infra/db/bun_helper"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/entity"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/utils/pagination"

	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
)

type BunRepairOrderRepository struct {
	db *bun.DB
}

func NewBunRepairOrderRepository(db *bun.DB) application.Repository {
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
	var model repairOrderModel
	err := r.db.NewSelect().
		Model(&model).
		Where("id = ?", id).
		Scan(ctx)

	if bun_helper.IgnoreNoRows(err) != nil {
		return nil, err
	}

	if model.ID == uuid.Nil {
		return nil, nil
	}

	return r.toEntityOrNil(&model), nil
}

func (r *BunRepairOrderRepository) List(ctx context.Context, pager pagination.Pagination) ([]*domain.RepairOrder, error) {
	var models []repairOrderModel

	query := r.db.NewSelect().
		Model(&models).
		Order("created_at DESC")

	if pager.Limit > 0 {
		query = query.Limit(pager.Limit)
	}
	if pager.Offset > 0 {
		query = query.Offset(pager.Offset)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return maps.Map(models, func(m repairOrderModel) *domain.RepairOrder {
		return r.toEntityOrNil(&m)
	}), nil
}

func (r *BunRepairOrderRepository) toEntityOrNil(m *repairOrderModel) *domain.RepairOrder {

	return &domain.RepairOrder{
		ID:         m.ID,
		CustomerID: m.CustomerID,
		VehicleID:  m.VehicleID,
		Status:     m.Status,
		Timestamps: entity.Timestamps{
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
	}
}

func (r *BunRepairOrderRepository) Save(ctx context.Context, repairOrder *domain.RepairOrder) error {
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

	return err
}

func (r *BunRepairOrderRepository) SaveIfApproved(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusApproved)
}

func (r *BunRepairOrderRepository) SaveIfInAwaitingApproval(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusAwaitingApproval)
}

func (r *BunRepairOrderRepository) SaveIfInDiagnostics(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusInDiagnostics)
}

func (r *BunRepairOrderRepository) SaveIfInExecution(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusInExecution)
}

func (r *BunRepairOrderRepository) SaveIfFinished(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusFinished)
}

func (r *BunRepairOrderRepository) saveIfStatus(
	ctx context.Context,
	ro *domain.RepairOrder,
	status repairorder.Status,
) error {
	_, err := r.db.NewUpdate().
		Model(ro).
		Where("id = ?", ro.ID).
		Where("status = ?", status).
		Exec(ctx)

	return err
}
