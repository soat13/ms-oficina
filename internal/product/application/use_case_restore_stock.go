package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/oficina-utils/pkg/maps"
)

type (
	RestoreStockInput struct {
		Products      map[uuid.UUID]int
		EstimateID    uuid.UUID
		RepairOrderID uuid.UUID
	}

	RestoreStock struct {
		repository Repository
	}
)

func NewRestoreStock(repository Repository) *RestoreStock {
	return &RestoreStock{
		repository: repository,
	}
}

func (uc *RestoreStock) Execute(ctx context.Context, in RestoreStockInput) error {
	productIDs := maps.Keys(in.Products)
	products, err := uc.repository.GetByIDs(ctx, productIDs)
	if err != nil {
		return err
	}

	for _, product := range products {
		quantity := in.Products[product.ID]
		newStock := product.Stock + quantity

		if err := product.UpdateStock(newStock); err != nil {
			return err
		}
	}

	return uc.repository.UpdateBatch(ctx, products)
}
