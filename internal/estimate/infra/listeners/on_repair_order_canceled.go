package listeners

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder/events"
)

func OnRepairOrderCanceled(cancel *application.Cancel, repository application.Repository) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var event events.Canceled
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		estimate, err := repository.GetByRepairOrderID(ctx, event.RepairOrderID)
		if err != nil {
			return err
		}

		if estimate == nil {
			return nil
		}

		return cancel.Execute(ctx, application.CancelInput{
			ID: estimate.ID,
		})
	}
}
