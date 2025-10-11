package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
	"github.com/soat13/fase-1-oficina/pkg/utils/pagination"
)

type (
	VehicleReader interface {
		Exists(ctx context.Context, id uuid.UUID) (bool, error)
	}

	CustomerReader interface {
		Exists(ctx context.Context, id uuid.UUID) (bool, error)
	}

	Repository interface {
		GetById(ctx context.Context, id uuid.UUID) (*domain.RepairOrder, error)
		List(ctx context.Context, pager pagination.Pagination) ([]*domain.RepairOrder, error)
		Create(ctx context.Context, repairOrder *domain.RepairOrder) error
		SaveIfApproved(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfInAwaitingApproval(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfInDiagnostics(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfInExecution(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfFinished(ctx context.Context, ro *domain.RepairOrder) error
	}
)
