package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
	"github.com/soat13/oficina-utils/pkg/pagination"
)

type (
	VehicleReader interface {
		Exists(ctx context.Context, id uuid.UUID) (bool, error)
	}

	CustomerReader interface {
		Exists(ctx context.Context, id uuid.UUID) (bool, error)
	}

	ServiceReader interface {
		ExistByIds(ctx context.Context, ids []uuid.UUID) (bool, error)
	}

	ProductReader interface {
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]ProductView, error)
	}

	ProductView struct {
		ID    uuid.UUID
		Stock int
	}

	Repository interface {
		GetById(ctx context.Context, id uuid.UUID) (*domain.RepairOrder, error)
		List(ctx context.Context, pager pagination.Pagination) ([]*domain.RepairOrder, error)
		Create(ctx context.Context, repairOrder *domain.RepairOrder) error
		SaveCancellation(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfReceived(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfInDiagnostics(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfDiagnosticsFinished(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfInAwaitingApproval(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfApproved(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfInExecution(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfFinished(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfPaymentCreated(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfPaymentProcessing(ctx context.Context, ro *domain.RepairOrder) error
		SaveIfPaymentSucceeded(ctx context.Context, ro *domain.RepairOrder) error
		GetAverageExecutionTime(ctx context.Context) (*float64, error)
	}

	EventPublisher interface {
		PublishRepairOrderDiagnosticsFinished(
			ctx context.Context,
			RepairOrderID uuid.UUID,
			products map[uuid.UUID]int,
			services map[uuid.UUID]int,
		) error
		PublishRepairOrderCanceled(ctx context.Context, RepairOrderID uuid.UUID) error
		PublishRepairOrderFinished(ctx context.Context, RepairOrderID uuid.UUID) error
	}

	MetricsPublisher interface {
		RecordRepairOrderPhaseDuration(phase string, minutes float64)
	}
)
