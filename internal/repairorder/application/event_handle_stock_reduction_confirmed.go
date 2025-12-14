package application

import (
	"context"

	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

type HandleStockReductionConfirmed struct {
	repository Repository
}

func NewHandleStockReductionConfirmed(repository Repository) HandleStockReductionConfirmed {
	return HandleStockReductionConfirmed{
		repository: repository,
	}
}

func (h *HandleStockReductionConfirmed) Execute(ctx context.Context, evt events.StockReductionConfirmed) error {
	repairOrder, err := h.repository.GetById(ctx, evt.RepairOrderID)
	if err != nil {
		return err
	}

	if err := repairOrder.Approve(); err != nil {
		return err
	}

	return h.repository.SaveIfInAwaitingApproval(ctx, repairOrder)
}
