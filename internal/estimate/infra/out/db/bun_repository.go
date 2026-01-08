package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/pkg/db/bun_helper"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type (
	estimateModel struct {
		bun.BaseModel `bun:"table:estimates"`

		ID            uuid.UUID `bun:"id,pk,type:uuid"`
		RepairOrderID uuid.UUID `bun:"repair_order_id,type:uuid,notnull"`
		Status        string    `bun:"status,notnull"`
		CreatedAt     time.Time `bun:"created_at,nullzero,default:now()"`
		UpdatedAt     time.Time `bun:"updated_at,nullzero,default:now()"`
	}

	estimateItemModel struct {
		bun.BaseModel `bun:"table:estimate_items"`

		ID         uuid.UUID `bun:"id,pk,type:uuid"`
		EstimateID uuid.UUID `bun:"estimate_id,type:uuid,notnull"`
		ItemID     uuid.UUID `bun:"item_id,type:uuid,notnull"`
		ItemName   string    `bun:"item_name,notnull"`
		ItemType   string    `bun:"item_type,notnull"`
		PriceCents int64     `bun:"price,notnull"`
		Quantity   int       `bun:"quantity,notnull"`
		CreatedAt  time.Time `bun:"created_at,nullzero,default:now()"`
		UpdatedAt  time.Time `bun:"updated_at,nullzero,default:now()"`
	}
)

type BunRepository struct {
	db *bun.DB
}

func NewBunRepository(db *bun.DB) application.Repository {
	return &BunRepository{db: db}
}

func (r *BunRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Estimate, error) {
	var model estimateModel
	if err := r.db.NewSelect().
		Model(&model).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, bun_helper.IgnoreNoRows(err)
	}

	var itemRows []estimateItemModel
	if err := r.db.NewSelect().
		Model(&itemRows).
		Where("estimate_id = ?", model.ID).
		Scan(ctx); err != nil {
		return nil, bun_helper.IgnoreNoRows(err)
	}

	return toEntity(model, itemRows)
}

func (r *BunRepository) GetByRepairOrderID(ctx context.Context, repairOrderID uuid.UUID) (*domain.Estimate, error) {
	var model estimateModel
	if err := r.db.NewSelect().
		Model(&model).
		Where("repair_order_id = ?", repairOrderID).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, bun_helper.IgnoreNoRows(err)
	}

	var itemRows []estimateItemModel
	if err := r.db.NewSelect().
		Model(&itemRows).
		Where("estimate_id = ?", model.ID).
		Scan(ctx); err != nil {
		return nil, bun_helper.IgnoreNoRows(err)
	}

	return toEntity(model, itemRows)
}

func (r *BunRepository) Save(ctx context.Context, estimate *domain.Estimate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	estimateModel := estimateModel{
		ID:            estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
		Status:        string(estimate.Status),
		UpdatedAt:     time.Now(),
	}

	if _, err := tx.NewInsert().
		Model(&estimateModel).
		On(`CONFLICT (id) DO UPDATE
		      SET status = EXCLUDED.status,
		          updated_at = EXCLUDED.updated_at`).
		Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewDelete().
		Model((*estimateItemModel)(nil)).
		Where("estimate_id = ?", estimateModel.ID).
		Exec(ctx); err != nil {
		return err
	}

	items := estimate.Items()
	if len(items) > 0 {
		itemRows := make([]estimateItemModel, 0, len(items))
		now := time.Now()
		for _, item := range items {
			itemRows = append(itemRows, estimateItemModel{
				ID:         uuid.New(),
				EstimateID: estimateModel.ID,
				ItemID:     item.ID,
				ItemName:   item.Name,
				ItemType:   string(item.Type),
				PriceCents: item.Price.Cents,
				Quantity:   item.Quantity,
				CreatedAt:  now,
				UpdatedAt:  now,
			})
		}
		if _, err := tx.NewInsert().Model(&itemRows).Exec(ctx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *BunRepository) SaveIfAwaitingStock(ctx context.Context, estimate *domain.Estimate) error {
	return r.saveIfStatus(ctx, estimate, domain.StatusAwaitingStock)
}

func (r *BunRepository) SaveIfAwaitingApproval(ctx context.Context, estimate *domain.Estimate) error {
	return r.saveIfStatus(ctx, estimate, domain.StatusAwaitingApproval)
}

func (r *BunRepository) saveIfStatus(ctx context.Context, estimate *domain.Estimate, status domain.Status) error {
	_, err := r.db.NewUpdate().
		Model(toModel(estimate)).
		WherePK().
		Where("status = ?", status).
		Exec(ctx)
	return err
}

func toModel(estimate *domain.Estimate) *estimateModel {
	return &estimateModel{
		ID:            estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
		Status:        string(estimate.Status),
		CreatedAt:     estimate.CreatedAt,
		UpdatedAt:     estimate.UpdatedAt,
	}
}

func toEntity(er estimateModel, itemRows []estimateItemModel) (*domain.Estimate, error) {
	estimate, err := domain.NewEstimate(er.RepairOrderID, er.CreatedAt, er.UpdatedAt)
	if err != nil {
		return nil, err
	}
	estimate.ID = er.ID
	estimate.Status = domain.Status(er.Status)

	for _, r := range itemRows {
		_ = estimate.AddItem(
			r.ItemID,
			r.ItemName,
			money.Money{Cents: r.PriceCents},
			r.Quantity,
			domain.ItemType(r.ItemType),
		)
	}
	return estimate, nil
}
