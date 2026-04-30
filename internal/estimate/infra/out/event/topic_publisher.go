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

type TopicPublisher struct {
	publisher messaging.TopicPublisher
}

func NewTopicPublisher(publisher messaging.TopicPublisher) application.TopicPublisher {
	return &TopicPublisher{publisher: publisher}
}

func (p *TopicPublisher) PublishApproved(ctx context.Context, estimate domain.Estimate) error {
	products := make(map[uuid.UUID]int, 0)
	for _, item := range estimate.Items() {
		if item.Type == domain.ProductItemType {
			products[item.ID] = item.Quantity
		}
	}

	event := estimateEvent.EstimateApproved{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
		Products:      products,
		TotalEstimate: estimate.Total(),
	}
	return p.publisher.Publish(ctx, messaging.TopicMessage{EventName: event.Topic(), Payload: event})
}
