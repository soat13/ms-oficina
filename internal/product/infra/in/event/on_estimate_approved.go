package event

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/product/application"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

func OnEstimateApproved(handle application.HandleEstimateApproved) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var eventPayload events.EstimateApproved
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			return err
		}

		return handle.Execute(ctx, eventPayload)
	}
}
