package db

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/auth/domain"
	"github.com/uptrace/bun"
)

type UserModel struct {
	bun.BaseModel `bun:"table:users"`

	ID           uuid.UUID `bun:",pk,type:uuid"`
	Name         string
	Email        string
	PasswordHash string
}

type BunUserRepository struct {
	db *bun.DB
}

func NewBunUserRepository(db *bun.DB) *BunUserRepository {
	return &BunUserRepository{db: db}
}

func (r *BunUserRepository) toDomain(m UserModel) domain.User {
	return domain.User{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
	}
}

func (r *BunUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	var m UserModel
	err := r.db.NewSelect().Model(&m).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(m), nil
}

func (r *BunUserRepository) Create(ctx context.Context, u domain.User) error {
	m := UserModel{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}
	_, err := r.db.NewInsert().Model(&m).Exec(ctx)
	return err
}

var _ interface {
	FindByEmail(context.Context, string) (domain.User, error)
	Create(context.Context, domain.User) error
} = (*BunUserRepository)(nil)
