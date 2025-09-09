package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	app "github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

func NewProductCatalogReader(db *bun.DB) app.ProductCatalogReader {
	return &catalogReader{db: db, table: "products"}
}

func NewServiceCatalogReader(db *bun.DB) app.ServiceCatalogReader {
	return &catalogReader{db: db, table: "services"}
}

type catalogReader struct {
	db    *bun.DB
	table string
}

type catalogRow struct {
	ID    uuid.UUID
	Name  string
	Price int64
}

func (r *catalogReader) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]app.CatalogItemView, error) {
	if len(ids) == 0 {
		return []app.CatalogItemView{}, nil
	}

	var rows []catalogRow
	if err := r.db.NewSelect().
		TableExpr("? AS c", bun.Ident(r.table)).
		ColumnExpr("c.id, c.name, c.price").
		Where("c.id IN (?)", bun.In(ids)).
		Scan(ctx, &rows); err != nil {
		return nil, err
	}

	out := make([]app.CatalogItemView, 0, len(rows))
	for _, row := range rows {
		m, err := money.New(row.Price)
		if err != nil {
			return nil, err
		}
		out = append(out, app.CatalogItemView{
			ID:    row.ID,
			Name:  row.Name,
			Price: m,
		})
	}
	return out, nil
}
