package listeners

import (
	"context"
	"encoding/json"

	"github.com/soat13/fase-1-oficina/internal/product/application"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/events"
)

func OnEstimateApproved(reduceStock *application.ReduceStock) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var eventPayload estimateEvent.EstimateApproved
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			return err
		}

		input := application.ReduceStockInput{
			Products:      eventPayload.Products,
			EstimateID:    eventPayload.EstimateID,
			RepairOrderID: eventPayload.RepairOrderID,
		}

		return reduceStock.Execute(ctx, input)
	}
}
