package application

import (
	"context"

	"github.com/google/uuid"
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
	}

	ProductCatalogReader interface {
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]CatalogItemView, error)
	}

	ServiceCatalogReader interface {
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]CatalogItemView, error)
	}
)
