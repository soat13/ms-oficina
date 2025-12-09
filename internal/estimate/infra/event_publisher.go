package infra

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/ports/event"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/events"
	eventPublisher "github.com/soat13/fase-1-oficina/pkg/event"
)

type (
	EventPublisher struct {
		bus event.Bus
	}
)

func NewEventPublisher(bus event.Bus) application.EventPublisher {
	return &EventPublisher{
		bus: bus,
	}
}

func (e *EventPublisher) PublishApproved(ctx context.Context, estimate domain.Estimate) error {
	products := make(map[uuid.UUID]int, 0)
	for _, item := range estimate.Items() {
		if item.Type == domain.ProductItemType {
			products[item.ID] = item.Quantity
		}
	}

	return eventPublisher.Publish(ctx, e.bus, estimateEvent.EstimateApproved{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
		Products:      products,
	})
}

func (e *EventPublisher) PublishRejected(ctx context.Context, estimate domain.Estimate) error {
	return eventPublisher.Publish(ctx, e.bus, estimateEvent.EstimateRejected{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
	})
}

func (e *EventPublisher) PublishCreated(ctx context.Context, estimate domain.Estimate) error {
	return eventPublisher.Publish(ctx, e.bus, estimateEvent.EstimateCreated{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
	})
}
