package listeners

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	estimateEvents "github.com/soat13/fase-1-oficina/internal/shared/events"
)

func OnEstimateRejected(cancel *application.Cancel) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var event estimateEvents.EstimateRejected
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		return cancel.Execute(ctx, application.CancelInput{RepairOrderID: event.RepairOrderID})
	}
}
