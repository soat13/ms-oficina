package application

import (
	"context"

	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

type HandleEstimateCreated struct {
	repository       Repository
	metricsPublisher MetricsPublisher
}

func NewHandleEstimateCreated(repository Repository, metricsPublisher MetricsPublisher) HandleEstimateCreated {
	return HandleEstimateCreated{
		repository:       repository,
		metricsPublisher: metricsPublisher,
	}
}

func (h *HandleEstimateCreated) Execute(ctx context.Context, evt events.EstimateCreated) error {
	repairOrder, err := h.repository.GetById(ctx, evt.RepairOrderID)

	if err != nil {
		return err
	}

	if err := repairOrder.MoveToAwaitingApproval(); err != nil {
		return err
	}

	if err := h.repository.SaveIfDiagnosticsFinished(ctx, repairOrder); err != nil {
		return err
	}

	return nil
}
