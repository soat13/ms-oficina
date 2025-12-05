package listeners

import (
	"context"
	"encoding/json"
	"time"

	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	repairOrderEvents "github.com/soat13/fase-1-oficina/internal/shared/events"
)

func OnDiagnosticsFinished(create *application.Create) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var event repairOrderEvents.RepairOrderDiagnosticsFinished
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		input := application.CreateInput{
			RepairOrderID: event.RepairOrderID,
			Products:      event.Products,
			Services:      event.Services,
			Now:           time.Now(),
		}

		return create.Execute(ctx, input)
	}
}
