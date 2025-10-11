package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/pkg/maps"
)

type (
	ReduceStockInput struct {
		Products map[uuid.UUID]int
	}

	ReduceStock struct {
		repo ProductRepository
	}
)

func NewReduceStock(repo ProductRepository) *ReduceStock {
	return &ReduceStock{repo: repo}
}

func (uc *ReduceStock) Execute(ctx context.Context, in ReduceStockInput) error {
	productIDs := maps.Keys(in.Products)
	products, err := uc.repo.GetByIDs(ctx, productIDs)
	if err != nil {
		return err
	}

	for _, product := range products {
		quantity := in.Products[product.ID]
		newStock := max(product.Stock-quantity, 0)

		if err := product.UpdateStock(newStock); err != nil {
			return err
		}
	}

	return uc.repo.UpdateBatch(ctx, products)
}
