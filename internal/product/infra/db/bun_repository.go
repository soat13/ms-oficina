package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	app "github.com/soat13/fase-1-oficina/internal/product/application"
	"github.com/soat13/fase-1-oficina/internal/product/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/infra/db/bun_helper"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/money"
	"github.com/soat13/fase-1-oficina/pkg/utils/pagination"
)

type productModel struct {
	bun.BaseModel `bun:"table:products"`

	ID        uuid.UUID `bun:",pk,type:uuid"`
	Name      string    `bun:",notnull"`
	Price     int64     `bun:",notnull"`
	Stock     int       `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,default:now()"`
	UpdatedAt time.Time `bun:",nullzero,default:now()"`
}

type BunRepository struct {
	db *bun.DB
}

func NewBunRepository(db *bun.DB) app.Repository {
	return &BunRepository{db: db}
}

func (r *BunRepository) Create(ctx context.Context, product *domain.Product) error {
	model := toModel(product)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return err
}

func (r *BunRepository) Update(ctx context.Context, product *domain.Product) error {
	model := toModel(product)
	_, err := r.db.NewUpdate().
		Model(model).
		Column("name", "price", "stock", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (r *BunRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&productModel{ID: id}).WherePK().Exec(ctx)
	return err
}

func (r *BunRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var model productModel
	err := r.db.NewSelect().Model(&model).Where("id = ?", id).Scan(ctx)

	if err := bun_helper.IgnoreNoRows(err); err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *BunRepository) List(ctx context.Context, pager pagination.Pagination) ([]*domain.Product, error) {
	var models []productModel
	if err := r.db.NewSelect().
		Model(&models).
		Order("name ASC").
		Limit(pager.Limit).
		Offset(pager.Offset).
		Scan(ctx); err != nil {
		return nil, err
	}

	return maps.MapPtr(models, toDomain), nil
}

func (r *BunRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
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

func toDomain(model *productModel) *domain.Product {
	if model == nil || model.ID == uuid.Nil {
		return nil
	}

	price, _ := money.New(model.Price)

	product, err := domain.NewProduct(model.ID, model.Name, price, model.Stock, model.CreatedAt, model.UpdatedAt)
	if err != nil {
		return nil
	}

	return product
}
