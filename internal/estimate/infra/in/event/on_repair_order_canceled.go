package event

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

func OnRepairOrderCanceled(handle application.HandleRepairOrderCanceled) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var event events.RepairOrderCanceled
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		return handle.Execute(ctx, event)
	}
}
