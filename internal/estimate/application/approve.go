package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
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

	return a.publishEvent(ctx, *estimate)
}

func (a *Approve) publishEvent(ctx context.Context, estimate domain.Estimate) error {
	products := make(map[uuid.UUID]int, 0)
	for _, item := range estimate.Items() {
		if item.Type == domain.ProductItemType {
			products[item.ID] = item.Quantity
		}
	}

	event := estimateEvent.StockReduceRequested{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		EstimateID:    estimate.ID,
		RepairOrderID: estimate.RepairOrderID,
		Products:      products,
	}

	b, _ := json.Marshal(event)
	return a.eventBus.Publish(ctx, event.Topic(), b)
}
