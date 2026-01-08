package db

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/service/application"
	"github.com/soat13/fase-1-oficina/internal/service/domain"
	"github.com/soat13/fase-1-oficina/pkg/db/bun_helper"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/money"
	"github.com/soat13/fase-1-oficina/pkg/pagination"
	"github.com/uptrace/bun"
)

type serviceModel struct {
	bun.BaseModel `bun:"table:services"`

	ID        uuid.UUID `bun:",pk,type:uuid"`
	Name      string    `bun:",notnull"`
	Price     int64     `bun:",notnull"`
	Currency  string    `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,default:now()"`
	UpdatedAt time.Time `bun:",nullzero,default:now()"`
}

type BunServiceRepository struct {
	db *bun.DB
}

func NewBunServiceRepository(db *bun.DB) application.Repository {
	return &BunServiceRepository{db: db}
}

func (repo *BunServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	model := toModel(service)
	_, err := repo.db.NewInsert().Model(model).Exec(ctx)
	return err
}

func (repo *BunServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	model := toModel(service)
	_, err := repo.db.NewUpdate().
		Model(model).
		Column("name", "price", "currency", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (repo *BunServiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.NewDelete().Model(&serviceModel{ID: id}).WherePK().Exec(ctx)
	return bun_helper.HandleDeleteError(err)
}

func (repo *BunServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	model := serviceModel{ID: id}
	err := repo.db.NewSelect().Model(&model).WherePK().Scan(ctx)

	if err := bun_helper.IgnoreNoRows(err); err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (repo *BunServiceRepository) List(ctx context.Context, pager pagination.Pagination) ([]*domain.Service, error) {
	var rows []serviceModel
	if err := repo.db.NewSelect().
		Model(&rows).
		Order("name ASC").
		Limit(pager.Limit).
		Offset(pager.Offset).
		Scan(ctx); err != nil {
		return nil, err
	}

	return maps.MapPtr(rows, toDomain), nil
}

func (repo *BunServiceRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	return repo.db.NewSelect().
		Model((*serviceModel)(nil)).
		Where("name = ?", strings.ToLower(name)).
		Exists(ctx)
}

func toModel(service *domain.Service) *serviceModel {
	return &serviceModel{
		ID:        service.ID,
		Name:      strings.ToLower(service.Name),
		Price:     service.Price.Cents,
		CreatedAt: service.CreatedAt,
		UpdatedAt: service.UpdatedAt,
	}
}

func toDomain(model *serviceModel) *domain.Service {
	if model == nil || model.ID == uuid.Nil {
		return nil
	}

	price, _ := money.New(model.Price)

	service, err := domain.NewService(model.ID, model.Name, price, model.CreatedAt, model.UpdatedAt)
	if err != nil {
		return nil
	}

	return service
}
