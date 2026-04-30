package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/repairorder/application"
	"github.com/uptrace/bun"
)

type VehicleReader struct {
	db *bun.DB
}

func NewVehicleReader(db *bun.DB) application.VehicleReader {
	return &VehicleReader{db: db}
}

type VehicleModel struct {
	bun.BaseModel `bun:"table:vehicles"`

	ID uuid.UUID `bun:"id,pk"`
}

func (r *VehicleReader) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.db.NewSelect().
		Model((*VehicleModel)(nil)).
		Where("id = ?", id).
		Exists(ctx)
}
