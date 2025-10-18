package application

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	productEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/product"

	"github.com/soat13/fase-1-oficina/pkg/maps"
)

type (
	ReduceStockInput struct {
		Products      map[uuid.UUID]int
		EstimateID    uuid.UUID
		RepairOrderID uuid.UUID
	}

	ReduceStock struct {
		repository ProductRepository
		eventBus   eventbus.Bus
	}
)

func NewReduceStock(repository ProductRepository, eventbus eventbus.Bus) *ReduceStock {
	return &ReduceStock{
		repository: repository,
		eventBus:   eventbus,
	}
}

func (uc *ReduceStock) Execute(ctx context.Context, in ReduceStockInput) error {
	if err := uc.updateStock(ctx, in); err != nil {
		if errors.Is(err, ErrInsufficientStock) {
			return uc.publishEvent(ctx, in)
		}
		return err
	}
	return nil
}

func (uc *ReduceStock) updateStock(ctx context.Context, in ReduceStockInput) error {
	productIDs := maps.Keys(in.Products)
	products, err := uc.repository.GetByIDs(ctx, productIDs)
	if err != nil {
		return err
	}

	for _, product := range products {
		quantity := in.Products[product.ID]
		newStock := product.Stock - quantity

		if newStock < 0 {
			return ErrInsufficientStock
		}

		if err := product.UpdateStock(newStock); err != nil {
			return err
		}
	}

	return uc.repository.UpdateBatch(ctx, products)
}

func (uc *ReduceStock) publishEvent(ctx context.Context, in ReduceStockInput) error {
	insufficientDetected := productEvent.StockInsufficientDetected{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    in.EstimateID,
		RepairOrderID: in.RepairOrderID,
	}

	payload, err := json.Marshal(insufficientDetected)
	if err != nil {
		return err
	}

	if err := uc.eventBus.Publish(ctx, insufficientDetected.Topic(), payload); err != nil {
		return err
	}

	return nil
}
