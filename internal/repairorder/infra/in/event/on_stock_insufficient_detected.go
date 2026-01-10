package event

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

func OnStockInsufficientDetected(handle application.HandleStockInsufficient) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var event events.StockInsufficientDetected
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		return handle.Execute(ctx, event)
	}
}
