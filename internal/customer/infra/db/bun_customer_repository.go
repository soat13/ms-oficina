package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/infra/db/bun_helper"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/document"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/phone"
	"github.com/uptrace/bun"

	app "github.com/soat13/fase-1-oficina/internal/customer/application"
	"github.com/soat13/fase-1-oficina/internal/customer/domain"
)

type customerModel struct {
	bun.BaseModel `bun:"table:customers"`

	ID           uuid.UUID `bun:",pk,type:uuid"`
	Name         string    `bun:",notnull"`
	Document     string    `bun:",notnull"`
	DocumentType string    `bun:",notnull"`
	Email        string    `bun:",notnull"`
	PhoneNumber  string    `bun:",notnull"`
	CreatedAt    time.Time `bun:",nullzero,default:now()"`
	UpdatedAt    time.Time `bun:",nullzero,default:now()"`
}

type BunCustomerRepository struct {
	db *bun.DB
}

func NewBunCustomerRepository(db *bun.DB) app.CustomerRepository {
	return &BunCustomerRepository{db: db}
}

func (repo *BunCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	model := toModel(customer)
	_, err := repo.db.NewInsert().Model(model).Exec(ctx)
	return err
}

func (repo *BunCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	model := toModel(customer)
	_, err := repo.db.NewUpdate().
		Model(model).
		Column("name", "document", "document_type", "email", "phone_number", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (repo *BunCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.NewDelete().Model(&customerModel{ID: id}).WherePK().Exec(ctx)
	return err
}

func (repo *BunCustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	model := customerModel{ID: id}
	err := repo.db.NewSelect().Model(&model).WherePK().Scan(ctx)

	if err := bun_helper.IgnoreNoRows(err); err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (repo *BunCustomerRepository) List(ctx context.Context, limit int, offset int) ([]*domain.Customer, error) {
	var rows []customerModel
	if err := repo.db.NewSelect().
		Model(&rows).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx); err != nil {
		return nil, err
	}

	out := maps.MapPtr(rows, toDomain)
	return out, nil
}

func (repo *BunCustomerRepository) ExistsByDocument(ctx context.Context, document string) (bool, error) {
	return repo.db.NewSelect().
		Model((*customerModel)(nil)).
		Where("document = ?", document).
		Exists(ctx)
}

func (repo *BunCustomerRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return repo.db.NewSelect().
		Model((*customerModel)(nil)).
		Where("email = ?", email).
		Exists(ctx)
}

func toModel(c *domain.Customer) *customerModel {
	return &customerModel{
		ID:           c.ID,
		Name:         c.Name,
		Document:     c.Document.Value,
		DocumentType: c.Document.TypeString(),
		PhoneNumber:  c.PhoneNumber.String(),
		Email:        c.Email.String(),
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func toDomain(model *customerModel) *domain.Customer {
	if model == nil || model.ID == uuid.Nil {
		return nil
	}

	document, _ := document.New(model.Document)
	phoneNumber, _ := phone.New(model.PhoneNumber)
	email, _ := email.New(model.Email)

	customer, err := domain.NewCustomer(model.ID, model.Name, document, phoneNumber, email, model.CreatedAt)
	if err != nil {
		return nil
	}

	customer.UpdatedAt = model.UpdatedAt
	return customer
}
