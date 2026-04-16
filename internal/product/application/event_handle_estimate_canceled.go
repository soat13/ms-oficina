package application

import (
	"context"

	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

type HandleEstimateCanceled struct {
	restoreStock *RestoreStock
}

func NewHandleEstimateCanceled(restoreStock *RestoreStock) HandleEstimateCanceled {
	return HandleEstimateCanceled{
		restoreStock: restoreStock,
	}
}

func (h *HandleEstimateCanceled) Execute(ctx context.Context, evt events.EstimateCanceled) error {
	if len(evt.Products) == 0 {
		return nil
	}

	return h.restoreStock.Execute(ctx, RestoreStockInput{
		Products:      evt.Products,
		EstimateID:    evt.EstimateID,
		RepairOrderID: evt.RepairOrderID,
	})
}
