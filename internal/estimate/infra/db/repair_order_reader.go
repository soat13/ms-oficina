package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/uptrace/bun"

	estimateApplication "github.com/soat13/fase-1-oficina/internal/estimate/application"
)

type RepairOrderReader struct {
	db *bun.DB
}

func NewRepairOrderReader(db *bun.DB) *RepairOrderReader {
	return &RepairOrderReader{db: db}
}

type repairOrderModel struct {
	bun.BaseModel `bun:"table:repair_orders"`

	ID     uuid.UUID `bun:"id,pk"`
	Status string    `bun:"status"`
}

type repairOrderItemModel struct {
	bun.BaseModel `bun:"table:repair_order_items"`

	RepairOrderID uuid.UUID `bun:"repair_order_id"`
	ItemID        uuid.UUID `bun:"item_id"`
	ItemType      string    `bun:"item_type"`
	Quantity      int       `bun:"quantity"`
}

func (r *RepairOrderReader) GetByID(ctx context.Context, id uuid.UUID) (*estimateApplication.RepairOrderView, error) {
	var m repairOrderModel

	if err := r.db.NewSelect().
		Model(&m).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	if m.ID == uuid.Nil {
		return nil, nil
	}

	var items []repairOrderItemModel
	if err := r.db.NewSelect().
		Model(&items).
		Where("repair_order_id = ?", id).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	view := &estimateApplication.RepairOrderView{
		ID:       m.ID,
		Status:   repairorder.Status(m.Status),
		Services: make(map[uuid.UUID]int),
		Products: make(map[uuid.UUID]int),
	}

	for _, item := range items {
		switch item.ItemType {
		case "service":
			view.Services[item.ItemID] = item.Quantity
		case "product":
			view.Products[item.ItemID] = item.Quantity
		}
	}

	return view, nil
}
