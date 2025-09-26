package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type (
	estimateModel struct {
		bun.BaseModel `bun:"table:estimates"`

		ID        uuid.UUID `bun:"id,pk,type:uuid"`
		RepairID  uuid.UUID `bun:"repair_id,type:uuid,notnull"`
		Status    string    `bun:"status,notnull"`
		CreatedAt time.Time `bun:"created_at,nullzero,default:now()"`
		UpdatedAt time.Time `bun:"updated_at,nullzero,default:now()"`
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

type BunEstimateRepository struct {
	db *bun.DB
}

func NewBunEstimateRepository(db *bun.DB) application.EstimateRepository {
	return &BunEstimateRepository{db: db}
}

func (r *BunEstimateRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Estimate, error) {
	var model estimateModel
	if err := r.db.NewSelect().
		Model(&model).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	}

	var itemRows []estimateItemModel
	if err := r.db.NewSelect().
		Model(&itemRows).
		Where("estimate_id = ?", model.ID).
		Scan(ctx); err != nil {
		return nil, err
	}

	return toEntity(model, itemRows)
}

func (r *BunEstimateRepository) Save(ctx context.Context, estimate *domain.Estimate) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	er := estimateModel{
		ID:       estimate.ID,
		RepairID: estimate.RepairOrderID,
		Status:   string(estimate.Status),
	}

	if _, err := tx.NewInsert().Model(&er).Exec(ctx); err != nil {
		return err
	}
	items := estimate.Items()
	itemRows := make([]estimateItemModel, 0, len(items))
	now := time.Now()
	for _, item := range items {
		itemRows = append(itemRows, estimateItemModel{
			ID:         uuid.New(),
			EstimateID: er.ID,
			ItemID:     item.ID,
			ItemName:   item.Name,
			ItemType:   string(item.Type),
			PriceCents: item.Price.Cents,
			Quantity:   item.Quantity,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}

	if len(itemRows) > 0 {
		if _, err := tx.NewInsert().Model(&itemRows).Exec(ctx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func toEntity(er estimateModel, itemRows []estimateItemModel) (*domain.Estimate, error) {
	estimate, err := domain.NewEstimate(er.RepairID, er.CreatedAt, er.UpdatedAt)
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
