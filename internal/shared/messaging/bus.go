package messaging

import "context"

type (
	Publisher interface {
		Publish(ctx context.Context, topic string, payload []byte) error
	}

	Handler func(ctx context.Context, topic string, payload []byte) error
)
