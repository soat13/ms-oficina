package listeners

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/shared/events/estimate"
)

func OnEstimateApproved(repository application.Repository) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var event estimate.Approved
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		repairOrder, err := repository.GetById(ctx, event.RepairOrderID)
		if err != nil {
			return err
		}

		if err := repairOrder.MoveToApproved(); err != nil {
			return err
		}

		_, err = repository.Save(ctx, repairOrder)
		return err
	}
}
