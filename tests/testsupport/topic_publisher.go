package testsupport

import (
	"context"

	"github.com/soat13/oficina-utils/pkg/messaging"
)

// fanoutTopicPublisher simulates SNS → SQS fanout for integration tests by
// forwarding a single topic publish to every SQS topic subscribed to it. The
// sync broker routes messages by EventName, so publishing to each queue topic
// triggers the matching in-memory handler.
type fanoutTopicPublisher struct {
	sender        messaging.QueueSender
	subscriptions map[string][]string
}

func NewFanoutTopicPublisher(sender messaging.QueueSender, subscriptions map[string][]string) messaging.TopicPublisher {
	return &fanoutTopicPublisher{sender: sender, subscriptions: subscriptions}
}

func (p *fanoutTopicPublisher) Publish(ctx context.Context, msg messaging.TopicMessage) error {
	for _, queueTopic := range p.subscriptions[msg.EventName] {
		if err := p.sender.Send(ctx, messaging.QueueMessage{
			EventName: queueTopic,
			Payload:   msg.Payload,
			GroupID:   msg.GroupID,
		}); err != nil {
			return err
		}
	}
	return nil
}
