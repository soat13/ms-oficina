package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/soat13/oficina-utils/pkg/maps"
)

type (
	ReduceStockInput struct {
		Products      map[uuid.UUID]int
		EstimateID    uuid.UUID
		RepairOrderID uuid.UUID
	}

	ReduceStock struct {
		repository     Repository
		eventPublisher EventPublisher
		topicPublisher TopicPublisher
	}
)

func NewReduceStock(repository Repository, eventPublisher EventPublisher, topicPublisher TopicPublisher) *ReduceStock {
	return &ReduceStock{
		repository:     repository,
		eventPublisher: eventPublisher,
		topicPublisher: topicPublisher,
	}
}

func (uc *ReduceStock) Execute(ctx context.Context, in ReduceStockInput) error {
	if err := uc.updateStock(ctx, in); err != nil {
		if !errors.Is(err, ErrInsufficientStock) {
			return err
		}

		return uc.eventPublisher.PublishStockInsufficientDetected(ctx, in.EstimateID, in.RepairOrderID)
	}

	return uc.topicPublisher.PublishStockReductionConfirmed(ctx, in.EstimateID, in.RepairOrderID)
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
