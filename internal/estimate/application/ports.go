package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/estimate/domain"
	"github.com/soat13/ms-oficina/internal/shared/repairorder"
	"github.com/soat13/oficina-utils/pkg/money"
)

type (
	RepairOrderView struct {
		ID     uuid.UUID
		Status repairorder.Status
	}

	CatalogItemView struct {
		ID    uuid.UUID
		Name  string
		Price money.Money
		Stock int
	}

	RepairOrderReader interface {
		GetByID(ctx context.Context, id uuid.UUID) (*RepairOrderView, error)
	}

	ProductCatalogReader interface {
		Exists(ctx context.Context, id uuid.UUID) (bool, error)
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]CatalogItemView, error)
	}

	ServiceCatalogReader interface {
		Exists(ctx context.Context, id uuid.UUID) (bool, error)
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]CatalogItemView, error)
	}

	EventPublisher interface {
		PublishRejected(ctx context.Context, estimate domain.Estimate) error
		PublishCreated(ctx context.Context, estimate domain.Estimate) error
		PublishCanceled(ctx context.Context, estimate domain.Estimate) error
	}

	TopicPublisher interface {
		PublishApproved(ctx context.Context, estimate domain.Estimate) error
	}

	Repository interface {
		GetByID(ctx context.Context, id uuid.UUID) (*domain.Estimate, error)
		GetByRepairOrderID(ctx context.Context, repairOrderID uuid.UUID) (*domain.Estimate, error)
		Save(ctx context.Context, estimate *domain.Estimate) error
		SaveIfAwaitingStock(ctx context.Context, estimate *domain.Estimate) error
		SaveIfAwaitingApproval(ctx context.Context, estimate *domain.Estimate) error
	}
)
