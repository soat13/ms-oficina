package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/pkg/db/bun_helper"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/entity"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/utils/pagination"

	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
)

type BunRepairOrderRepository struct {
	db *bun.DB
}

type repairOrderModel struct {
	bun.BaseModel `bun:"table:repair_orders"`

	ID                   uuid.UUID `bun:"id,pk,type:uuid"`
	CustomerID           uuid.UUID `bun:"customer_id,type:uuid,notnull"`
	VehicleID            uuid.UUID `bun:"vehicle_id,type:uuid,notnull"`
	Status               string    `bun:"status,notnull"`
	ExecutionTimeMinutes *int64    `bun:"execution_time_minutes"`
	CreatedAt            time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt            time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

var _ bun.BeforeAppendModelHook = (*repairOrderModel)(nil)

func (m *repairOrderModel) BeforeAppendModel(ctx context.Context, q bun.Query) error {
	switch q.(type) {
	case *bun.InsertQuery:
		if m.CreatedAt.IsZero() {
			now := time.Now()
			m.CreatedAt = now
			m.UpdatedAt = now
		}
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}

func NewBunRepairOrderRepository(db *bun.DB) application.Repository {
	return &BunRepairOrderRepository{db: db}
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
	timestamps := entity.NewTimestamps(m.CreatedAt, m.UpdatedAt)

	return &domain.RepairOrder{
		ID:                   m.ID,
		CustomerID:           m.CustomerID,
		VehicleID:            m.VehicleID,
		Status:               repairorder.Status(m.Status),
		ExecutionTimeMinutes: m.ExecutionTimeMinutes,
		Timestamps:           &timestamps,
	}
}

func (r *BunRepairOrderRepository) Create(ctx context.Context, repairOrder *domain.RepairOrder) error {
	_, err := r.db.NewInsert().
		Model(toModel(repairOrder)).
		On("CONFLICT (id) DO UPDATE").
		Set("customer_id = EXCLUDED.customer_id").
		Set("vehicle_id = EXCLUDED.vehicle_id").
		Set("status = EXCLUDED.status").
		Set("execution_time_minutes = EXCLUDED.execution_time_minutes").
		Exec(ctx)
	return err
}

func (r *BunRepairOrderRepository) SaveCancellation(ctx context.Context, ro *domain.RepairOrder) error {
	deniedStatusList := []repairorder.Status{
		repairorder.StatusReleased,
		repairorder.StatusCanceled,
	}

	_, err := r.db.NewUpdate().
		Model(toModel(ro)).
		WherePK().
		Where("status not in (?)", bun.In(deniedStatusList)).
		Exec(ctx)

	return err
}

func (r *BunRepairOrderRepository) SaveIfReceived(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusReceived)
}

func (r *BunRepairOrderRepository) SaveIfInDiagnostics(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusInDiagnostics)
}

func (r *BunRepairOrderRepository) SaveIfDiagnosticsFinished(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusDiagnosticsFinished)
}

func (r *BunRepairOrderRepository) SaveIfInAwaitingApproval(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusAwaitingApproval)
}

func (r *BunRepairOrderRepository) SaveIfApproved(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusApproved)
}

func (r *BunRepairOrderRepository) SaveIfInExecution(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusInExecution)
}

func (r *BunRepairOrderRepository) SaveIfFinished(ctx context.Context, ro *domain.RepairOrder) error {
	return r.saveIfStatus(ctx, ro, repairorder.StatusFinished)
}

func (r *BunRepairOrderRepository) GetAverageExecutionTime(ctx context.Context) (*float64, error) {
	var avg float64

	validStatusList := []repairorder.Status{
		repairorder.StatusReleased,
		repairorder.StatusFinished,
	}

	err := r.db.NewSelect().
		Model((*repairOrderModel)(nil)).
		ColumnExpr("AVG(execution_time_minutes) as average").
		Where("execution_time_minutes IS NOT NULL").
		Where("status in (?)", bun.In(validStatusList)).
		Scan(ctx, &avg)

	if err != nil {
		return nil, err
	}

	return &avg, nil
}

func (r *BunRepairOrderRepository) saveIfStatus(ctx context.Context, ro *domain.RepairOrder, status repairorder.Status) error {
	_, err := r.db.NewUpdate().
		Model(toModel(ro)).
		WherePK().
		Where("status = ?", status).
		Exec(ctx)
	return err
}

func toModel(ro *domain.RepairOrder) *repairOrderModel {
	return &repairOrderModel{
		ID:                   ro.ID,
		CustomerID:           ro.CustomerID,
		VehicleID:            ro.VehicleID,
		Status:               string(ro.Status),
		ExecutionTimeMinutes: ro.ExecutionTimeMinutes,
		CreatedAt:            ro.CreatedAt,
		UpdatedAt:            ro.UpdatedAt,
	}
}
