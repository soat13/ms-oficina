package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

type EventPublisher struct {
	publisher messaging.Publisher
}

func NewEventPublisher(publisher messaging.Publisher) application.EventPublisher {
	return &EventPublisher{publisher: publisher}
}

func (e *EventPublisher) PublishRepairOrderDiagnosticsFinished(
	ctx context.Context,
	RepairOrderID uuid.UUID,
	products map[uuid.UUID]int,
	services map[uuid.UUID]int,
) error {
	orderDiagnosticsFinished := estimateEvent.RepairOrderDiagnosticsFinished{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		RepairOrderID: RepairOrderID,
		Products:      products,
		Services:      services,
	}

	return messaging.Publish(ctx, e.publisher, orderDiagnosticsFinished)
}

func (e *EventPublisher) PublishRepairOrderCanceled(ctx context.Context, RepairOrderID uuid.UUID) error {
	evt := estimateEvent.RepairOrderCanceled{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		RepairOrderID: RepairOrderID,
	}

	return messaging.Publish(ctx, e.publisher, evt)
}

func (e *EventPublisher) PublishRepairOrderFinished(ctx context.Context, RepairOrderID uuid.UUID) error {
	evt := estimateEvent.RepairOrderFinished{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		RepairOrderID: RepairOrderID,
	}

	return messaging.Publish(ctx, e.publisher, evt)
}
