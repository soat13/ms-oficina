package messaging

import "context"

type (
	Bus interface {
		Publish(ctx context.Context, topic string, payload []byte) error
		Subscribe(topic string, h Handler)
	}

	Handler func(ctx context.Context, topic string, payload []byte) error
)
