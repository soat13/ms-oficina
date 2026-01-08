package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/uptrace/bun"
)

type ServiceReader struct {
	db *bun.DB
}

func NewServiceReader(db *bun.DB) application.ServiceReader {
	return &ServiceReader{db: db}
}

type ServiceModel struct {
	bun.BaseModel `bun:"table:services"`

	ID uuid.UUID `bun:"id,pk"`
}

func (r *ServiceReader) ExistByIds(ctx context.Context, ids []uuid.UUID) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}

	count, err := r.db.NewSelect().
		Model((*ServiceModel)(nil)).
		Where("id IN (?)", bun.In(ids)).
		Count(ctx)

	if err != nil {
		return false, err
	}
	return count == len(ids), nil
}
