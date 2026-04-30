package event

import (
	"context"

	"github.com/soat13/ms-oficina/internal/estimate/application"
	repairOrderEvents "github.com/soat13/ms-oficina/internal/shared/events"
	"github.com/soat13/oficina-utils/pkg/messaging"
)

func OnDiagnosticsFinished(handle application.HandleDiagnosticsFinished) messaging.Handler {
	return func(ctx context.Context, msg messaging.Message) error {
		event, err := messaging.DecodePayload[repairOrderEvents.RepairOrderDiagnosticsFinished](msg)
		if err != nil {
			return err
		}

		return handle.Execute(ctx, *event)
	}
}
