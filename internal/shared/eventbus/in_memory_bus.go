package eventbus

import (
	"context"
	"sync"
)

type Handler func(ctx context.Context, topic string, payload []byte) error

type inMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewInMemoryBus() Bus { return &inMemoryBus{handlers: map[string][]Handler{}} }

func (b *inMemoryBus) Publish(ctx context.Context, topic string, payload []byte) error {
	b.mu.RLock()
	hs := append([]Handler{}, b.handlers[topic]...)
	b.mu.RUnlock()
	for _, h := range hs {
		if err := h(ctx, topic, payload); err != nil {
			return err
		}
	}
	return nil
}

func (b *inMemoryBus) Subscribe(topic string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], h)
}
