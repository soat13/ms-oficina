package listeners

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	productEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/product"
)

func OnStockReduceConfirmed(repository application.Repository) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var event productEvent.StockReduceConfirmed
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		repairOrder, err := repository.GetById(ctx, event.RepairOrderID)
		if err != nil {
			return err
		}

		if err := repairOrder.Approve(); err != nil {
			return err
		}

		return repository.SaveIfInAwaitingApproval(ctx, repairOrder)
	}
}
