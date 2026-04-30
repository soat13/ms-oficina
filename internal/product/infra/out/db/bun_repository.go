package db

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/product/application"
	"github.com/soat13/ms-oficina/internal/product/domain"
	"github.com/soat13/oficina-utils/pkg/db/bun_helper"
	"github.com/soat13/oficina-utils/pkg/entity"
	"github.com/soat13/oficina-utils/pkg/maps"
	"github.com/soat13/oficina-utils/pkg/money"
	"github.com/soat13/oficina-utils/pkg/pagination"
	"github.com/uptrace/bun"
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

type BunProductRepository struct {
	db *bun.DB
}

func NewBunProductRepository(db *bun.DB) application.Repository {
	return &BunProductRepository{db: db}
}

func (r *BunProductRepository) Create(ctx context.Context, product *domain.Product) error {
	model := toModel(product)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return err
}

func (r *BunProductRepository) Update(ctx context.Context, product *domain.Product) error {
	model := toModel(product)
	_, err := r.db.NewUpdate().
		Model(model).
		Column("name", "price", "stock", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (r *BunProductRepository) UpdateBatch(ctx context.Context, products []*domain.Product) error {
	if len(products) == 0 {
		return nil
	}

	models := maps.Map(products, toModel)
	_, err := r.db.NewUpdate().
		Model(&models).
		Column("name", "price", "stock", "updated_at").
		Bulk().
		Exec(ctx)
	return err
}

func (r *BunProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&productModel{ID: id}).WherePK().Exec(ctx)
	return bun_helper.HandleDeleteError(err)
}

func (r *BunProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	model := productModel{ID: id}
	err := r.db.NewSelect().Model(&model).WherePK().Scan(ctx)

	if err := bun_helper.IgnoreNoRows(err); err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *BunProductRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error) {
	if len(ids) == 0 {
		return []*domain.Product{}, nil
	}

	var models []productModel
	err := r.db.NewSelect().
		Model(&models).
		Where("id IN (?)", bun.In(ids)).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return maps.MapPtr(models, toDomain), nil
}

func (r *BunProductRepository) List(ctx context.Context, pager pagination.Pagination) ([]*domain.Product, error) {
	var rows []productModel
	if err := r.db.NewSelect().
		Model(&rows).
		Order("name ASC").
		Limit(pager.Limit).
		Offset(pager.Offset).
		Scan(ctx); err != nil {
		return nil, err
	}

	return maps.MapPtr(rows, toDomain), nil
}

func (r *BunProductRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	return r.db.NewSelect().
		Model((*productModel)(nil)).
		Where("name = ?", strings.ToLower(name)).
		Exists(ctx)
}

func toModel(p *domain.Product) *productModel {
	return &productModel{
		ID:        p.ID,
		Name:      strings.ToLower(p.Name),
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

	product, err := domain.NewProduct(model.ID, model.Name, price, model.Stock)
	if err != nil {
		return nil
	}

	product.Timestamps = entity.NewTimestamps(model.CreatedAt, model.UpdatedAt)
	return product
}
