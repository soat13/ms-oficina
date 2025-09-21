package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	app "github.com/soat13/fase-1-oficina/internal/customer/application"
	"github.com/soat13/fase-1-oficina/internal/customer/domain"
)

type customerModel struct {
	bun.BaseModel `bun:"table:customers"`

	ID           uuid.UUID `bun:",pk,type:uuid"`
	Name         string    `bun:",notnull"`
	Cellphone    string    `bun:",notnull"`
	Document     string    `bun:",notnull"`
	DocumentType string    `bun:",notnull"`
	CreatedAt    time.Time `bun:",nullzero,default:now()"`
	UpdatedAt    time.Time `bun:",nullzero,default:now()"`
}

type BunCustomerRepository struct {
	db *bun.DB
}

func NewBunCustomerRepository(db *bun.DB) *BunCustomerRepository {
	return &BunCustomerRepository{db: db}
}

func (r *BunCustomerRepository) Create(ctx context.Context, c *domain.Customer) error {
	m := toModel(c)
	_, err := r.db.NewInsert().Model(m).Exec(ctx)
	return err
}

func (r *BunCustomerRepository) Update(ctx context.Context, c *domain.Customer) error {
	m := toModel(c)
	_, err := r.db.NewUpdate().
		Model(m).
		Column("name", "cellphone", "document", "document_type", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (r *BunCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&customerModel{ID: id}).WherePK().Exec(ctx)
	return err
}

func (r *BunCustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	var m customerModel
	err := r.db.NewSelect().Model(&m).Where("id = ?", id).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, app.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *BunCustomerRepository) List(ctx context.Context, limit, offset int) ([]*domain.Customer, error) {
	var rows []customerModel
	if err := r.db.NewSelect().
		Model(&rows).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx); err != nil {
		return nil, err
	}
	out := make([]*domain.Customer, 0, len(rows))
	for i := range rows {
		out = append(out, toDomain(&rows[i]))
	}
	return out, nil
}

func (r *BunCustomerRepository) ExistsByDocument(ctx context.Context, document string) (bool, error) {
	return r.db.NewSelect().
		Model((*customerModel)(nil)).
		Where("document = ?", document).
		Exists(ctx)
}

func toModel(c *domain.Customer) *customerModel {
	return &customerModel{
		ID:           c.ID,
		Name:         c.Name,
		Cellphone:    c.Cellphone,
		Document:     c.Document,
		DocumentType: c.DocumentType,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func toDomain(m *customerModel) *domain.Customer {
	return &domain.Customer{
		ID:           m.ID,
		Name:         m.Name,
		Cellphone:    m.Cellphone,
		Document:     m.Document,
		DocumentType: m.DocumentType,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
