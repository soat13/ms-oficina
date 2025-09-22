package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/events/estimate"
)

type (
	ApproveInput struct {
		ID uuid.UUID
	}

	Approve struct {
		repository Repository
		eventBus   eventbus.Bus
	}
)

func NewApproveEstimate(repository Repository, eventBus eventbus.Bus) *Approve {
	return &Approve{
		repository: repository,
		eventBus:   eventBus,
	}
}

func (a *Approve) Execute(ctx context.Context, input ApproveInput) error {
	estimate, err := a.repository.GetByID(ctx, input.ID)
	if err != nil {
		return ErrEstimateNotFound
	}

	if err := estimate.MoveToAwaitingStock(); err != nil {
		return err
	}

	if err := a.repository.Save(ctx, estimate); err != nil {
		return err
	}

	if err := a.publishEvent(ctx, *estimate); err != nil {
		return err
	}

	return nil
}

func (a *Approve) publishEvent(ctx context.Context, estimate domain.Estimate) error {
	event := estimateEvent.ApprovedByCustomer{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
	}

	b, _ := json.Marshal(event)
	if err := a.eventBus.Publish(ctx, event.Topic(), b); err != nil {
		return err
	}

	return nil
}
