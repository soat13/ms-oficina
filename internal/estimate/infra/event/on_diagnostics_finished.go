package event

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	repairOrderEvents "github.com/soat13/fase-1-oficina/internal/shared/events"
)

func OnDiagnosticsFinished(handle application.HandleDiagnosticsFinished) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var event repairOrderEvents.RepairOrderDiagnosticsFinished
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		return handle.Execute(ctx, event)
	}
}
