package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/auth/application"
	"github.com/soat13/fase-1-oficina/internal/shared/authz"
	"github.com/soat13/fase-1-oficina/pkg/db/bun_helper"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/password"
	"github.com/uptrace/bun"
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

type UserReader struct {
	db *bun.DB
}

func NewBunUserReader(db *bun.DB) application.UserReader {
	return &UserReader{db: db}
}

func (repo *UserReader) GetByCPF(ctx context.Context, email string) (*application.UserView, error) {
	var model userModel
	err := repo.db.NewSelect().
		Model(&model).
		Where("document = ?", email).
		Scan(ctx)

	if err := bun_helper.IgnoreNoRows(err); err != nil {
		return nil, err
	}

	roles, _ := authz.NewRoles(model.Roles)
	pwd, _ := password.FromHash(model.Password)

	return &application.UserView{
		ID:       model.ID,
		Name:     model.Name,
		Email:    model.Email,
		Password: pwd,
		Roles:    roles,
	}, nil
}
