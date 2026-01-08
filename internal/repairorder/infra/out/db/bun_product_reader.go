package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/uptrace/bun"
)

type ProductReader struct {
	db *bun.DB
}

func NewProductReader(db *bun.DB) application.ProductReader {
	return &ProductReader{db: db}
}

type ProductModel struct {
	bun.BaseModel `bun:"table:products"`

	ID    uuid.UUID `bun:"id,pk"`
	Stock int       `bun:"stock"`
}

func (r *ProductReader) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]application.ProductView, error) {
	if len(ids) == 0 {
		return []application.ProductView{}, nil
	}

	var models []ProductModel
	err := r.db.NewSelect().
		Model(&models).
		Column("id", "stock").
		Where("id IN (?)", bun.In(ids)).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	views := make([]application.ProductView, len(models))
	for i, model := range models {
		views[i] = application.ProductView{
			ID:    model.ID,
			Stock: model.Stock,
		}
	}

	return views, nil
}
