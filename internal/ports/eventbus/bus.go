package eventbus

import (
	"context"
)

type Handler func(ctx context.Context, topic string, payload []byte) error

type Bus interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	Subscribe(topic string, h Handler)
}
