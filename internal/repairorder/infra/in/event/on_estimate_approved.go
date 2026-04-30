package event

import (
	"context"

	"github.com/soat13/ms-oficina/internal/repairorder/application"
	"github.com/soat13/ms-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

func OnEstimateApproved(handle application.HandleEstimateApproved) messaging.Handler {
	return func(ctx context.Context, msg messaging.Message) error {
		event, err := messaging.DecodePayload[events.EstimateApproved](msg)
		if err != nil {
			return err
		}

		return handle.Execute(ctx, *event)
	}
}
