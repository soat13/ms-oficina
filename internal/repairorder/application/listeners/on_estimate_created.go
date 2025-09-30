package listeners

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/shared/events/estimate"
)

func OnEstimateCreated(repository application.Repository) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var event estimate.Created
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		repairOrder, err := repository.GetById(ctx, event.RepairOrderID)
		if err != nil {
			return err
		}

		if err := repairOrder.MoveToAwaitingApproval(); err != nil {
			return err
		}

		return repository.SaveIfInDiagnostics(ctx, repairOrder)
	}
}
