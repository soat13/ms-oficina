package event_test

import (
	"context"
	"testing"

	"github.com/soat13/fase-1-oficina/pkg/event"
	"github.com/stretchr/testify/require"
)

type badEvent struct {
	C chan int `json:"c"`
}

func (e badEvent) Topic() string { return "bad.topic" }

func TestPublish(t *testing.T) {
	t.Run("should returns error when marshal fails", func(t *testing.T) {
		ctx := context.Background()
		invalidEvent := badEvent{C: make(chan int)}

		err := event.Publish(ctx, nil, invalidEvent)
		require.Error(t, err)
	})
}
