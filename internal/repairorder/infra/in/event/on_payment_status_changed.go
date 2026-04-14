package event

import (
	"context"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

func OnPaymentStatusChanged(handle application.HandlePaymentStatusChanged) messaging.Handler {
	return func(ctx context.Context, msg messaging.Message) error {
		event, err := messaging.DecodePayload[events.PaymentStatusChanged](msg)
		if err != nil {
			return err
		}

		return handle.Execute(ctx, *event)
	}
}
