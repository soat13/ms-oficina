package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/product/application"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

type EventPublisher struct {
	publisher messaging.QueueSender
}

func NewEventPublisher(publisher messaging.QueueSender) application.EventPublisher {
	return &EventPublisher{publisher: publisher}
}

func (e *EventPublisher) PublishStockInsufficientDetected(
	ctx context.Context,
	estimateID uuid.UUID,
	repairOrderID uuid.UUID,
) error {
	event := estimateEvent.StockInsufficientDetected{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimateID,
		RepairOrderID: repairOrderID,
	}
	return e.publisher.Send(ctx, messaging.QueueMessage{EventName: event.Topic(), Payload: event})
}

func (e *EventPublisher) PublishStockReduceConfirmed(
	ctx context.Context,
	estimateID uuid.UUID,
	repairOrderID uuid.UUID,
) error {
	event := estimateEvent.StockReductionConfirmed{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimateID,
		RepairOrderID: repairOrderID,
	}
	return e.publisher.Send(ctx, messaging.QueueMessage{EventName: event.Topic(), Payload: event})
}
