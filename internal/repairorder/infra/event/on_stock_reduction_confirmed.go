package event

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

func OnStockReduceConfirmed(handle application.HandleStockReductionConfirmed) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var evt events.StockReductionConfirmed
		if err := json.Unmarshal(payload, &evt); err != nil {
			return err
		}

		return handle.Execute(ctx, evt)
	}
}
