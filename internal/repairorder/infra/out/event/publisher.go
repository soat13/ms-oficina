package event

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
	"github.com/soat13/oficina-utils/pkg/money"
)

type EventPublisher struct {
	publisher messaging.QueueSender
}

func NewEventPublisher(publisher messaging.QueueSender) application.EventPublisher {
	return &EventPublisher{publisher: publisher}
}

func (e *EventPublisher) PublishRepairOrderDiagnosticsFinished(
	ctx context.Context,
	RepairOrderID uuid.UUID,
	products map[uuid.UUID]int,
	services map[uuid.UUID]int,
) error {
	event := estimateEvent.RepairOrderDiagnosticsFinished{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		RepairOrderID: RepairOrderID,
		Products:      products,
		Services:      services,
	}
	return e.publisher.Send(ctx, messaging.QueueMessage{EventName: event.Topic(), Payload: event})
}

func (e *EventPublisher) PublishRepairOrderCanceled(ctx context.Context, RepairOrderID uuid.UUID) error {
	event := estimateEvent.RepairOrderCanceled{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		RepairOrderID: RepairOrderID,
	}
	return e.publisher.Send(ctx, messaging.QueueMessage{EventName: event.Topic(), Payload: event})
}

func (e *EventPublisher) PublishPaymentRequest(ctx context.Context, ro *domain.RepairOrder) error {
	event := estimateEvent.PaymentRequest{
		RepairOrderID: ro.ID,
		Amount:        money.Money{Cents: ro.TotalEstimate.Cents},
		Description:   fmt.Sprintf("Ordem de serviço %s", ro.ID),
	}
	return e.publisher.Send(ctx, messaging.QueueMessage{EventName: event.Topic(), Payload: event})
}
