package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/ports/eventbus"
	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder/events"
)

type CancelInput struct {
	RepairOrderID uuid.UUID
}

type Cancel struct {
	Repository Repository
	Bus        eventbus.Bus
}

func NewCancel(repository Repository, bus eventbus.Bus) *Cancel {
	return &Cancel{
		Repository: repository,
		Bus:        bus,
	}
}

func (uc *Cancel) Execute(ctx context.Context, input CancelInput) error {
	repairorder, err := uc.Repository.GetById(ctx, input.RepairOrderID)
	if err != nil {
		return err
	}

	if repairorder == nil {
		return sharedRepairOrder.ErrRepairOrderNotFound
	}

	if repairorder.IsCancelled() {
		return nil
	}

	if err := repairorder.Cancel(); err != nil {
		return err
	}

	if err := uc.Repository.SaveCancellation(ctx, repairorder); err != nil {
		return err
	}

	return uc.publishEvent(ctx, *repairorder)
}

func (uc *Cancel) publishEvent(ctx context.Context, repairOrder domain.RepairOrder) error {
	event := events.Canceled{
		EventID:       uuid.New(),
		OccurredAt:    time.Now(),
		RepairOrderID: repairOrder.ID,
	}

	b, err := json.Marshal(event)

	if err != nil {
		return err
	}

	return uc.Bus.Publish(ctx, event.Topic(), b)
}
