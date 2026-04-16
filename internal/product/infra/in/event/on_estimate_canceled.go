package event

import (
	"context"

	"github.com/soat13/fase-1-oficina/internal/product/application"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

func OnEstimateCanceled(handle application.HandleEstimateCanceled) messaging.Handler {
	return func(ctx context.Context, msg messaging.Message) error {
		event, err := messaging.DecodePayload[events.EstimateCanceled](msg)
		if err != nil {
			return err
		}

		return handle.Execute(ctx, *event)
	}
}
