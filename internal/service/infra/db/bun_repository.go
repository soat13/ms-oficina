package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	app "github.com/soat13/fase-1-oficina/internal/service/application"
	"github.com/soat13/fase-1-oficina/internal/service/domain"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type serviceModel struct {
	bun.BaseModel `bun:"table:services"`

	ID        uuid.UUID `bun:",pk,type:uuid"`
	Name      string    `bun:",notnull"`
	Price     int64     `bun:",notnull"` // cents
	Currency  string    `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,default:now()"`
	UpdatedAt time.Time `bun:",nullzero,default:now()"`
}

type BunServiceRepository struct {
	db *bun.DB
}

func NewBunServiceRepository(db *bun.DB) *BunServiceRepository {
	return &BunServiceRepository{db: db}
}

func (r *BunServiceRepository) Create(ctx context.Context, s *domain.Service) error {
	m := toModel(s)
	_, err := r.db.NewInsert().Model(m).Exec(ctx)
	return err
}

func (r *BunServiceRepository) Update(ctx context.Context, s *domain.Service) error {
	m := toModel(s)
	_, err := r.db.NewUpdate().
		Model(m).
		Column("name", "price", "currency", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (r *BunServiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&serviceModel{ID: id}).WherePK().Exec(ctx)
	return err
}

func (r *BunServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	var m serviceModel
	err := r.db.NewSelect().Model(&m).Where("id = ?", id).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, app.ErrServiceNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *BunServiceRepository) List(ctx context.Context, limit, offset int) ([]*domain.Service, error) {
	var rows []serviceModel
	if err := r.db.NewSelect().
		Model(&rows).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]*domain.Service, 0, len(rows))
	for i := range rows {
		out = append(out, toDomain(&rows[i]))
	}
	return out, nil
}

func (r *BunServiceRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	return r.db.NewSelect().
		Model((*serviceModel)(nil)).
		Where("LOWER(name) = LOWER(?)", name).
		Exists(ctx)
}

func toModel(s *domain.Service) *serviceModel {
	return &serviceModel{
		ID:        s.ID,
		Name:      s.Name,
		Price:     s.Price.Cents,
		Currency:  s.Currency,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func toDomain(m *serviceModel) *domain.Service {
	price, _ := money.New(m.Price)

	return &domain.Service{
		ID:        m.ID,
		Name:      m.Name,
		Price:     price,
		Currency:  m.Currency,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
