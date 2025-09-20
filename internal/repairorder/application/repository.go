package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
)

type (
	Repository interface {
		GetById(ctx context.Context, id uuid.UUID) (*domain.RepairOrder, error)
		Save(ctx context.Context, repairOrder *domain.RepairOrder) (*domain.RepairOrder, error)
	}
)
