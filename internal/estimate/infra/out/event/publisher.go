package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/estimate/application"
	"github.com/soat13/ms-oficina/internal/estimate/domain"
	estimateEvent "github.com/soat13/ms-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

type EventPublisher struct {
	publisher messaging.QueueSender
}

func NewEventPublisher(publisher messaging.QueueSender) application.EventPublisher {
	return &EventPublisher{publisher: publisher}
}

func (e *EventPublisher) PublishRejected(ctx context.Context, estimate domain.Estimate) error {
	event := estimateEvent.EstimateRejected{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
	}
	return e.publisher.Send(ctx, messaging.QueueMessage{EventName: event.Topic(), Payload: event})
}

func (e *EventPublisher) PublishCreated(ctx context.Context, estimate domain.Estimate) error {
	event := estimateEvent.EstimateCreated{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
	}
	return e.publisher.Send(ctx, messaging.QueueMessage{EventName: event.Topic(), Payload: event})
}

func (e *EventPublisher) PublishCanceled(ctx context.Context, estimate domain.Estimate, wasApproved bool) error {
	products := make(map[uuid.UUID]int)
	if wasApproved {
		for _, item := range estimate.Items() {
			if item.Type == domain.ProductItemType {
				products[item.ID] = item.Quantity
			}
		}
	}

	event := estimateEvent.EstimateCanceled{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
		Products:      products,
	}
	return e.publisher.Send(ctx, messaging.QueueMessage{EventName: event.Topic(), Payload: event})
}
