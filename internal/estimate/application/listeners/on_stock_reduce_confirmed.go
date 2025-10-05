package listeners

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
	productEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/product"
)

func OnStockReduceConfirmed(repository application.Repository, eventBus eventbus.Bus) func(ctx context.Context, _ string, payload []byte) error {
	return func(ctx context.Context, _ string, payload []byte) error {
		var event productEvent.StockReduceConfirmed
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		estimate, err := repository.GetByID(ctx, event.EstimateID)
		if err != nil {
			return err
		}

		if err := estimate.Approve(); err != nil {
			return err
		}

		if err := repository.SaveIfAwaitingStock(ctx, estimate); err != nil {
			return err
		}
		return publishEvent(ctx, estimate, eventBus)
	}
}

func publishEvent(ctx context.Context, estimate *domain.Estimate, eventBus eventbus.Bus) error {
	event := estimateEvent.Approved{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
	}

	payload, _ := json.Marshal(event)
	if err := eventBus.Publish(ctx, event.Topic(), payload); err != nil {
		return err
	}

	return nil
}
