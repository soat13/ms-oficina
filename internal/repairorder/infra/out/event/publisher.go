package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
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

	return messaging.Publish(ctx, e.bus, orderDiagnosticsFinished)
}

func (e *EventPublisher) PublishRepairOrderCanceled(ctx context.Context, RepairOrderID uuid.UUID) error {
	orderDiagnosticsFinished := estimateEvent.RepairOrderCanceled{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		RepairOrderID: RepairOrderID,
	}

	return messaging.Publish(ctx, e.bus, orderDiagnosticsFinished)
}
