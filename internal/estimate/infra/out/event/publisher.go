package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

type EventPublisher struct {
	publisher messaging.Publisher
}

func NewEventPublisher(publisher messaging.Publisher) application.EventPublisher {
	return &EventPublisher{publisher: publisher}
}

func (e *EventPublisher) PublishApproved(ctx context.Context, estimate domain.Estimate) error {
	products := make(map[uuid.UUID]int, 0)
	for _, item := range estimate.Items() {
		if item.Type == domain.ProductItemType {
			products[item.ID] = item.Quantity
		}
	}

	return messaging.Publish(ctx, e.publisher, estimateEvent.EstimateApproved{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
		Products:      products,
	})
}

func (e *EventPublisher) PublishRejected(ctx context.Context, estimate domain.Estimate) error {
	return messaging.Publish(ctx, e.publisher, estimateEvent.EstimateRejected{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
	})
}

func (e *EventPublisher) PublishCreated(ctx context.Context, estimate domain.Estimate) error {
	return messaging.Publish(ctx, e.publisher, estimateEvent.EstimateCreated{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
	})
}

func (e *EventPublisher) PublishCanceled(ctx context.Context, estimate domain.Estimate) error {
	products := make(map[uuid.UUID]int)
	for _, item := range estimate.Items() {
		if item.Type == domain.ProductItemType {
			products[item.ID] = item.Quantity
		}
	}

	return messaging.Publish(ctx, e.publisher, estimateEvent.EstimateCanceled{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
		Products:      products,
	})
}
