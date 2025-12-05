package application

import (
	"context"

	"github.com/google/uuid"
	estimateDomain "github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type (
	RepairOrderStatus string

	RepairOrderView struct {
		ID     uuid.UUID
		Status repairorder.Status
	}

	RepairOrderReader interface {
		GetByID(ctx context.Context, id uuid.UUID) (*RepairOrderView, error)
	}

	CatalogItemView struct {
		ID    uuid.UUID
		Name  string
		Price money.Money
		Stock int
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
		PublishApproved(ctx context.Context, estimate estimateDomain.Estimate) error
		PublishRejected(ctx context.Context, estimate estimateDomain.Estimate) error
		PublishCreated(ctx context.Context, estimate estimateDomain.Estimate) error
	}
)
