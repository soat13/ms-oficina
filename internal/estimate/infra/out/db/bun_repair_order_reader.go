package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/repairorder"
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

	view := &estimateApplication.RepairOrderView{
		ID:     m.ID,
		Status: repairorder.Status(m.Status),
	}

	return view, nil
}
