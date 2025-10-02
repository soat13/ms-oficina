package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	DeleteInput struct {
		ID uuid.UUID
	}

	DeleteUser struct {
		repo UserRepository
	}
)

func NewDeleteUser(repo UserRepository) *DeleteUser {
	return &DeleteUser{repo: repo}
}

func (uc *DeleteUser) Execute(ctx context.Context, in DeleteInput) error {
	return uc.repo.Delete(ctx, in.ID)
}
