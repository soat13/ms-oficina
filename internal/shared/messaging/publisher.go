package messaging

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

func Publish(ctx context.Context, pub Publisher, event events.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return pub.Publish(ctx, event.Topic(), payload)
}
