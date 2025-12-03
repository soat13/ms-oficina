package listeners

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/ports/eventbus"
	"github.com/soat13/fase-1-oficina/internal/product/application"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
	productEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/product"
)

func OnEstimateApproved(reduceStock *application.ReduceStock, eventBus eventbus.Bus) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var event estimateEvent.Approved
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		input := application.ReduceStockInput{
			Products:      event.Products,
			EstimateID:    event.EstimateID,
			RepairOrderID: event.RepairOrderID,
		}
		if err := reduceStock.Execute(ctx, input); err != nil {
			return err
		}

		return publishEvent(ctx, event, eventBus)
	}
}

func publishEvent(ctx context.Context, event estimateEvent.Approved, eventBus eventbus.Bus) error {
	confirmedEvent := productEvent.StockReduceConfirmed{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    event.EstimateID,
		RepairOrderID: event.RepairOrderID,
	}

	payload, _ := json.Marshal(confirmedEvent)
	if err := eventBus.Publish(ctx, confirmedEvent.Topic(), payload); err != nil {
		return err
	}

	return nil
}
