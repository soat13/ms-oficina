package bus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/soat13/fase-1-oficina/pkg/event/bus"
	"github.com/stretchr/testify/require"
)

func TestInMemoryBus_Publish(t *testing.T) {
	t.Run("returns handler error when a handler fails", func(t *testing.T) {
		ctx := context.Background()
		b := bus.NewInMemoryBus()

		handlerErr := errors.New("handler failed")

		b.Subscribe("topic.test", func(ctx context.Context, topic string, payload []byte) error {
			return handlerErr
		})

		err := b.Publish(ctx, "topic.test", []byte(`{"x":1}`))
		require.ErrorIs(t, err, handlerErr)
	})

	t.Run("calls handlers in order and stops on first error", func(t *testing.T) {
		ctx := context.Background()
		b := bus.NewInMemoryBus()

		calls := 0
		handlerErr := errors.New("boom")

		b.Subscribe("topic.test", func(ctx context.Context, topic string, payload []byte) error {
			calls++
			return nil
		})
		b.Subscribe("topic.test", func(ctx context.Context, topic string, payload []byte) error {
			calls++
			return handlerErr
		})
		b.Subscribe("topic.test", func(ctx context.Context, topic string, payload []byte) error {
			calls++
			return nil
		})

		err := b.Publish(ctx, "topic.test", []byte("payload"))
		require.ErrorIs(t, err, handlerErr)
		require.Equal(t, 2, calls)
	})

	t.Run("returns nil when there are no handlers for topic", func(t *testing.T) {
		ctx := context.Background()
		b := bus.NewInMemoryBus()

		err := b.Publish(ctx, "no.handlers", []byte("payload"))
		require.NoError(t, err)
	})
}
