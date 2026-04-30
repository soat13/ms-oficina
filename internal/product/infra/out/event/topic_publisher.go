package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/product/application"
	productEvent "github.com/soat13/ms-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

type TopicPublisher struct {
	publisher messaging.TopicPublisher
}

func NewTopicPublisher(publisher messaging.TopicPublisher) application.TopicPublisher {
	return &TopicPublisher{publisher: publisher}
}

func (p *TopicPublisher) PublishStockReductionConfirmed(
	ctx context.Context,
	estimateID uuid.UUID,
	repairOrderID uuid.UUID,
) error {
	event := productEvent.StockReductionConfirmed{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimateID,
		RepairOrderID: repairOrderID,
	}

	return p.publisher.Publish(ctx, messaging.TopicMessage{EventName: event.Topic(), Payload: event})
}
