package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/product/application"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/events"
	"github.com/soat13/fase-1-oficina/internal/shared/messaging"
)

type EventPublisher struct {
	bus messaging.Bus
}

func NewEventPublisher(bus messaging.Bus) application.EventPublisher {
	return &EventPublisher{
		bus: bus,
	}
}

func (e *EventPublisher) PublishStockInsufficientDetected(
	ctx context.Context,
	estimateID uuid.UUID,
	repairOrderID uuid.UUID,
) error {
	insufficientDetected := estimateEvent.StockInsufficientDetected{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimateID,
		RepairOrderID: repairOrderID,
	}

	return messaging.Publish(ctx, e.bus, insufficientDetected)
}

func (e *EventPublisher) PublishStockReduceConfirmed(
	ctx context.Context,
	estimateID uuid.UUID,
	repairOrderID uuid.UUID,
) error {
	stockReduceConfirmed := estimateEvent.StockReductionConfirmed{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimateID,
		RepairOrderID: repairOrderID,
	}

	return messaging.Publish(ctx, e.bus, stockReduceConfirmed)
}
