package application

import (
	"context"

	"github.com/soat13/ms-oficina/internal/shared/events"
)

type HandleEstimateApproved struct {
	reduceStock *ReduceStock
}

func NewHandleEstimateApproved(reduceStock *ReduceStock) HandleEstimateApproved {
	return HandleEstimateApproved{
		reduceStock: reduceStock,
	}
}

func (h *HandleEstimateApproved) Execute(ctx context.Context, evt events.EstimateApproved) error {
	input := ReduceStockInput{
		Products:      evt.Products,
		EstimateID:    evt.EstimateID,
		RepairOrderID: evt.RepairOrderID,
	}

	return h.reduceStock.Execute(ctx, input)
}
