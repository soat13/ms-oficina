package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/infra/db/bun_helper"
	roleVO "github.com/soat13/fase-1-oficina/internal/user/domain/role"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/utils/pagination"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/document"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/password"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/phone"
	"github.com/uptrace/bun"

	app "github.com/soat13/fase-1-oficina/internal/user/application"
	"github.com/soat13/fase-1-oficina/internal/user/domain"
)

type userModel struct {
	bun.BaseModel `bun:"table:users"`

	ID           uuid.UUID `bun:",pk,type:uuid"`
	Name         string    `bun:",notnull"`
	Document     string    `bun:",notnull,unique"`
	DocumentType string    `bun:",notnull"`
	Email        string    `bun:",notnull,unique"`
	PhoneNumber  string    `bun:",notnull"`
	Password     string    `bun:",notnull"`
	Roles        []string  `bun:",array,notnull"`
	CreatedAt    time.Time `bun:",nullzero,default:now()"`
	UpdatedAt    time.Time `bun:",nullzero,default:now()"`
}

type BunUserRepository struct {
	db *bun.DB
}

func NewBunUserRepository(db *bun.DB) app.UserRepository {
	return &BunUserRepository{db: db}
}

func (repo *BunUserRepository) Create(ctx context.Context, user *domain.User) error {
	model := toModel(user)
	_, err := repo.db.NewInsert().Model(model).Exec(ctx)
	return err
}

func (repo *BunUserRepository) Update(ctx context.Context, user *domain.User) error {
	model := toModel(user)
	_, err := repo.db.NewUpdate().
		Model(model).
		Column("name", "document", "document_type", "email", "phone_number", "password", "roles", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (repo *BunUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.NewDelete().Model(&userModel{ID: id}).WherePK().Exec(ctx)
	return err
}

func (repo *BunUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	model := userModel{ID: id}
	err := repo.db.NewSelect().Model(&model).WherePK().Scan(ctx)

	if err := bun_helper.IgnoreNoRows(err); err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (repo *BunUserRepository) List(ctx context.Context, pager pagination.Pagination) ([]*domain.User, error) {
	var rows []userModel
	if err := repo.db.NewSelect().
		Model(&rows).
		Order("name ASC").
		Limit(pager.Limit).
		Offset(pager.Offset).
		Scan(ctx); err != nil {
		return nil, err
	}

	out := maps.MapPtr(rows, toDomain)
	return out, nil
}

func (repo *BunUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return repo.db.NewSelect().
		Model((*userModel)(nil)).
		Where("email = ?", email).
		Exists(ctx)
}

func (repo *BunUserRepository) ExistsByDocument(ctx context.Context, document string) (bool, error) {
	return repo.db.NewSelect().
		Model((*userModel)(nil)).
		Where("document = ?", document).
		Exists(ctx)
}

func toModel(u *domain.User) *userModel {
	return &userModel{
		ID:           u.ID,
		Name:         u.Name,
		Document:     u.Document.Value,
		DocumentType: u.Document.TypeString(),
		Email:        u.Email.String(),
		PhoneNumber:  u.PhoneNumber.String(),
		Password:     u.Password.Hash,
		Roles:        u.Roles.Strings(),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func toDomain(model *userModel) *domain.User {
	if model == nil || model.ID == uuid.Nil {
		return nil
	}

	doc, _ := document.New(model.Document)
	phoneNumber, _ := phone.New(model.PhoneNumber)
	em, _ := email.New(model.Email)
	pwd, _ := password.FromHash(model.Password)
	roles, _ := roleVO.NewRoles(model.Roles)

	user, err := domain.NewUser(model.ID, model.Name, doc, phoneNumber, em, pwd, roles, model.CreatedAt)
	if err != nil {
		return nil
	}

	user.UpdatedAt = model.UpdatedAt
	return user
}
