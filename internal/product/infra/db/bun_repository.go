package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	app "github.com/soat13/fase-1-oficina/internal/product/application"
	"github.com/soat13/fase-1-oficina/internal/product/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type productModel struct {
	bun.BaseModel `bun:"table:products"`

	ID        uuid.UUID `bun:",pk,type:uuid"`
	Name      string    `bun:",notnull"`
	Price     int64     `bun:",notnull"` // cents
	Stock     int       `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,default:now()"`
	UpdatedAt time.Time `bun:",nullzero,default:now()"`
}

type BunProductRepository struct {
	db *bun.DB
}

func NewBunProductRepository(db *bun.DB) *BunProductRepository {
	return &BunProductRepository{db: db}
}

func (r *BunProductRepository) Create(ctx context.Context, p *domain.Product) error {
	m := toModel(p)
	_, err := r.db.NewInsert().Model(m).Exec(ctx)
	return err
}

func (r *BunProductRepository) Update(ctx context.Context, p *domain.Product) error {
	m := toModel(p)
	_, err := r.db.NewUpdate().
		Model(m).
		Column("name", "price", "stock", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (r *BunProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&productModel{ID: id}).WherePK().Exec(ctx)
	return err
}

func (r *BunProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var m productModel
	err := r.db.NewSelect().Model(&m).Where("id = ?", id).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, app.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *BunProductRepository) List(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	var rows []productModel
	if err := r.db.NewSelect().
		Model(&rows).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]*domain.Product, 0, len(rows))
	for i := range rows {
		out = append(out, toDomain(&rows[i]))
	}
	return out, nil
}

func (r *BunProductRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	return r.db.NewSelect().
		Model((*productModel)(nil)).
		Where("LOWER(name) = LOWER(?)", name).
		Exists(ctx)
}

func toModel(p *domain.Product) *productModel {
	return &productModel{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price.Cents,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func toDomain(m *productModel) *domain.Product {
	price, _ := money.New(m.Price)
	return &domain.Product{
		ID:        m.ID,
		Name:      m.Name,
		Price:     price,
		Stock:     m.Stock,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
