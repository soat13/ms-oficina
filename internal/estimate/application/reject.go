package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/ports/eventbus"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
)

type (
	RejectInput struct {
		RepairOrderID uuid.UUID
	}

	Reject struct {
		repository Repository
		eventBus   eventbus.Bus
	}
)

func NewRejectEstimate(repository Repository, eventBus eventbus.Bus) *Reject {
	return &Reject{
		repository: repository,
		eventBus:   eventBus,
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

	return a.publishEvent(ctx, *estimate)
}

func (a *Reject) publishEvent(ctx context.Context, estimate domain.Estimate) error {
	event := estimateEvent.Rejected{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
	}

	b, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return a.eventBus.Publish(ctx, event.Topic(), b)
}
