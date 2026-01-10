package messaging

import (
	"context"
	"sync"

	"github.com/soat13/fase-1-oficina/internal/shared/messaging"
)

type inMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]messaging.Handler
}

func NewInMemoryBus() messaging.Bus { return &inMemoryBus{handlers: map[string][]messaging.Handler{}} }

func (b *inMemoryBus) Publish(ctx context.Context, topic string, payload []byte) error {
	b.mu.RLock()
	hs := append([]messaging.Handler{}, b.handlers[topic]...)
	b.mu.RUnlock()
	for _, h := range hs {
		if err := h(ctx, topic, payload); err != nil {
			return err
		}
	}
	return nil
}

func (b *inMemoryBus) Subscribe(topic string, h messaging.Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], h)
}
