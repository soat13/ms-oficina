package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	RejectInput struct {
		RepairOrderID uuid.UUID
	}

	Reject struct {
		repository     Repository
		eventPublisher EventPublisher
	}
)

func NewRejectEstimate(repository Repository, eventPublisher EventPublisher) *Reject {
	return &Reject{
		repository:     repository,
		eventPublisher: eventPublisher,
	}
}

func (a *Reject) Execute(ctx context.Context, input RejectInput) error {
	estimate, err := a.repository.GetByRepairOrderID(ctx, input.RepairOrderID)
	if err != nil {
		return err
	}

	if estimate == nil {
		return ErrEstimateNotFound
	}

	if err := estimate.Reject(); err != nil {
		return err
	}

	if err := a.repository.Save(ctx, estimate); err != nil {
		return err
	}

	return a.eventPublisher.PublishRejected(ctx, *estimate)
}
