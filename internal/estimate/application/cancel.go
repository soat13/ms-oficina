package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
)

type (
	CancelInput struct {
		ID uuid.UUID
	}

	Cancel struct {
		repository Repository
		eventBus   eventbus.Bus
	}
)

func NewCancelEstimate(repository Repository, eventBus eventbus.Bus) *Cancel {
	return &Cancel{
		repository: repository,
		eventBus:   eventBus,
	}
}

func (a *Cancel) Execute(ctx context.Context, input CancelInput) error {
	estimate, err := a.repository.GetByID(ctx, input.ID)
	if err != nil {
		return ErrEstimateNotFound
	}

	if estimate == nil {
		return ErrEstimateNotFound
	}

	if estimate.IsCanceled() || estimate.IsRejected() {
		return nil
	}

	if err := estimate.Cancel(); err != nil {
		return err
	}

	return a.repository.Save(ctx, estimate)
}
