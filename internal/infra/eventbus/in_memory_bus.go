package eventbus

import (
	"context"
	"sync"

	"github.com/soat13/fase-1-oficina/internal/ports/eventbus"
)

type inMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]eventbus.Handler
}

func NewInMemoryBus() eventbus.Bus { return &inMemoryBus{handlers: map[string][]eventbus.Handler{}} }

func (b *inMemoryBus) Publish(ctx context.Context, topic string, payload []byte) error {
	b.mu.RLock()
	hs := append([]eventbus.Handler{}, b.handlers[topic]...)
	b.mu.RUnlock()
	for _, h := range hs {
		if err := h(ctx, topic, payload); err != nil {
			return err
		}
	}
	return nil
}

func (b *inMemoryBus) Subscribe(topic string, h eventbus.Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], h)
}
