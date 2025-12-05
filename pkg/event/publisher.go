package event

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/ports/event"
)

func Publish(ctx context.Context, bus event.Bus, event event.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return bus.Publish(ctx, event.Topic(), payload)
}
