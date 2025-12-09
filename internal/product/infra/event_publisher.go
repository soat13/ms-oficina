package infra

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/ports/event"
	"github.com/soat13/fase-1-oficina/internal/product/application"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/events"
	eventPublisher "github.com/soat13/fase-1-oficina/pkg/event"
)

type EventPublisher struct {
	bus event.Bus
}

func NewEventPublisher(bus event.Bus) application.EventPublisher {
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

	return eventPublisher.Publish(ctx, e.bus, insufficientDetected)
}

func (e *EventPublisher) PublishStockReduceConfirmed(
	ctx context.Context,
	estimateID uuid.UUID,
	repairOrderID uuid.UUID,
) error {
	stockReduceConfirmed := estimateEvent.StockReduceConfirmed{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimateID,
		RepairOrderID: repairOrderID,
	}

	return eventPublisher.Publish(ctx, e.bus, stockReduceConfirmed)
}
