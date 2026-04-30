package application

import (
	"context"

	"github.com/soat13/ms-oficina/internal/shared/events"
)

type HandleRepairOrderCanceled struct {
	cancel     Cancel
	repository Repository
}

func NewHandleRepairOrderCanceled(cancel Cancel, repository Repository) HandleRepairOrderCanceled {
	return HandleRepairOrderCanceled{
		cancel:     cancel,
		repository: repository,
	}
}

func (h *HandleRepairOrderCanceled) Execute(ctx context.Context, evt events.RepairOrderCanceled) error {
	estimate, err := h.repository.GetByRepairOrderID(ctx, evt.RepairOrderID)
	if err != nil {
		return err
	}

	if estimate == nil {
		return nil
	}

	return h.cancel.Execute(ctx, CancelInput{
		ID: estimate.ID,
	})
}
