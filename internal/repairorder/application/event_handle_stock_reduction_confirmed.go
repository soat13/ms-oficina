package application

import (
	"context"
	"time"

	"github.com/soat13/ms-oficina/internal/shared/events"
)

type HandleStockReductionConfirmed struct {
	repository       Repository
	metricsPublisher MetricsPublisher
}

func NewHandleStockReductionConfirmed(repository Repository, metricsPublisher MetricsPublisher) HandleStockReductionConfirmed {
	return HandleStockReductionConfirmed{
		repository:       repository,
		metricsPublisher: metricsPublisher,
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

	if err := h.repository.ApproveIfAwaitingApproval(ctx, repairOrder); err != nil {
		return err
	}

	if repairOrder.Timestamps != nil {
		h.metricsPublisher.RecordRepairOrderPhaseDuration("awaiting_approval", time.Since(repairOrder.UpdatedAt).Minutes())
	}

	return nil
}
