package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
)

type (
	Repository interface {
		GetByID(ctx context.Context, id uuid.UUID) (*domain.Estimate, error)
		Save(ctx context.Context, estimate *domain.Estimate) error
	}
)
