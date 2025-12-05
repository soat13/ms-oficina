package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	ApproveInput struct {
		RepairOrderID uuid.UUID
	}

	Approve struct {
		repository     Repository
		eventPublisher EventPublisher
	}
)

func NewApproveEstimate(repository Repository, eventPublisher EventPublisher) *Approve {
	return &Approve{
		repository:     repository,
		eventPublisher: eventPublisher,
	}
}

func (a *Approve) Execute(ctx context.Context, input ApproveInput) error {
	estimate, err := a.repository.GetByRepairOrderID(ctx, input.RepairOrderID)
	if err != nil {
		return nil
	}

	if estimate == nil {
		return ErrEstimateNotFound
	}

	if err := estimate.Approve(); err != nil {
		return err
	}

	if err := a.repository.SaveIfAwaitingApproval(ctx, estimate); err != nil {
		return err
	}

	return a.eventPublisher.PublishApproved(ctx, *estimate)
}
