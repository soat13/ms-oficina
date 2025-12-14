package application

import (
	"context"

	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

type HandleEstimateRejected struct {
	cancel *Cancel
}

func NewHandleEstimateRejected(cancel *Cancel) HandleEstimateRejected {
	return HandleEstimateRejected{
		cancel: cancel,
	}
}

func (h *HandleEstimateRejected) Execute(ctx context.Context, evt events.EstimateRejected) error {
	return h.cancel.Execute(ctx, CancelInput{RepairOrderID: evt.RepairOrderID})
}
