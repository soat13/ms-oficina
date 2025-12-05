package event

import (
	"context"
)

type (
	Event interface {
		Topic() string
	}

	Bus interface {
		Publish(ctx context.Context, topic string, payload []byte) error
		Subscribe(topic string, h Handler)
	}

	Handler func(ctx context.Context, topic string, payload []byte) error
)
