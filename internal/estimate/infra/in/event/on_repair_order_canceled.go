package event

import (
	"context"

	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

func OnRepairOrderCanceled(handle application.HandleRepairOrderCanceled) messaging.Handler {
	return func(ctx context.Context, msg messaging.Message) error {
		event, err := messaging.DecodePayload[events.RepairOrderCanceled](msg)
		if err != nil {
			return err
		}

		return handle.Execute(ctx, *event)
	}
}
