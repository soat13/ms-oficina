package event

import (
	"context"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

func OnEstimateRejected(handler application.HandleEstimateRejected) messaging.Handler {
	return func(ctx context.Context, msg messaging.Message) error {
		event, err := messaging.DecodePayload[events.EstimateRejected](msg)
		if err != nil {
			return err
		}

		return handler.Execute(ctx, *event)
	}
}
