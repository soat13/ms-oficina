package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/repairorder/application"
	"github.com/uptrace/bun"
)

type CustomerReader struct {
	db *bun.DB
}

func NewCustomerReader(db *bun.DB) application.CustomerReader {
	return &CustomerReader{db: db}
}

type CustomerModel struct {
	bun.BaseModel `bun:"table:customers"`

	ID uuid.UUID `bun:"id,pk"`
}

func (r *CustomerReader) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.db.NewSelect().
		Model((*CustomerModel)(nil)).
		Where("id = ?", id).
		Exists(ctx)
}
