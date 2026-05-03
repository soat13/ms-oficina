package application

import (
	"context"

	"github.com/google/uuid"
)

type (
	CancelInput struct {
		ID uuid.UUID
	}

	Cancel struct {
		repository     Repository
		eventPublisher EventPublisher
	}
)

func NewCancelEstimate(repository Repository, eventPublisher EventPublisher) *Cancel {
	return &Cancel{
		repository:     repository,
		eventPublisher: eventPublisher,
	}
}

func (a *Cancel) Execute(ctx context.Context, input CancelInput) error {
	estimate, err := a.repository.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}
	if estimate == nil {
		return ErrEstimateNotFound
	}

	if estimate.IsCanceled() || estimate.IsRejected() {
		return nil
	}

	shouldPublishEvent := estimate.IsApproved() || estimate.IsAwaitingStock()

	if err := estimate.Cancel(); err != nil {
		return err
	}

	if err := a.repository.Save(ctx, estimate); err != nil {
		return err
	}

	if shouldPublishEvent {
		return a.eventPublisher.PublishCanceled(ctx, *estimate)
	}

	return nil
}
