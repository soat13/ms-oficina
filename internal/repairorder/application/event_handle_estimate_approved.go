package application

import (
	"context"

	"github.com/soat13/ms-oficina/internal/shared/events"
	sharedRepairOrder "github.com/soat13/ms-oficina/internal/shared/repairorder"
)

type HandleEstimateApproved struct {
	repository Repository
}

func NewHandleEstimateApproved(repository Repository) HandleEstimateApproved {
	return HandleEstimateApproved{
		repository: repository,
	}
}

func (h *HandleEstimateApproved) Execute(ctx context.Context, evt events.EstimateApproved) error {
	repairOrder, err := h.repository.GetById(ctx, evt.RepairOrderID)
	if err != nil {
		return err
	}

	if repairOrder == nil {
		return sharedRepairOrder.ErrRepairOrderNotFound
	}

	repairOrder.UpdateTotalEstimate(&evt.TotalEstimate)

	return h.repository.SaveTotalEstimateIfAwaitingApproval(ctx, repairOrder)
}
