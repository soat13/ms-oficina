package application

import (
	"context"

	"github.com/google/uuid"
)

type DeleteInput struct {
	ID uuid.UUID
}

type DeleteService struct {
	repo ServiceRepository
}

func NewDeleteService(repo ServiceRepository) *DeleteService {
	return &DeleteService{repo: repo}
}

func (uc *DeleteService) Execute(ctx context.Context, in DeleteInput) error {
	return uc.repo.Delete(ctx, in.ID)
}
