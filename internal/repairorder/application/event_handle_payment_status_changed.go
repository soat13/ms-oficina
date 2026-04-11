package application

import (
	"context"
	"time"

	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusSucceeded Status = "SUCCEEDED"
	StatusFailed    Status = "FAILED"
	StatusError     Status = "ERROR"
)

type HandlePaymentStatusChanged struct {
	repository       Repository
	metricsPublisher MetricsPublisher
}

func NewHandlePaymentStatusChanged(repository Repository, metricsPublisher MetricsPublisher) HandlePaymentStatusChanged {
	return HandlePaymentStatusChanged{
		repository:       repository,
		metricsPublisher: metricsPublisher,
	}
}

func (h *HandlePaymentStatusChanged) Execute(ctx context.Context, evt events.PaymentStatusChanged) error {
	repairOrder, err := h.repository.GetById(ctx, evt.ID)
	if err != nil {
		return err
	}

	if repairOrder == nil {
		return sharedRepairOrder.ErrRepairOrderNotFound
	}

	if evt.PaymentURL != nil {
		repairOrder.UpdatePaymentURL(evt.PaymentURL)
	}

	return h.handleStatusChanged(ctx, repairOrder, evt.Status)
}

func (h *HandlePaymentStatusChanged) handleStatusChanged(ctx context.Context, repairOrder *domain.RepairOrder, status string) error {
	switch Status(status) {
	case StatusPending:
		if err := repairOrder.PaymentCreated(); err != nil {
			return err
		}
		if err := h.repository.SaveIfFinished(ctx, repairOrder); err != nil {
			return err
		}
		if repairOrder.Timestamps != nil {
			h.metricsPublisher.RecordRepairOrderPhaseDuration("finished", time.Since(repairOrder.UpdatedAt).Minutes())
		}
	case StatusSucceeded:
		if err := repairOrder.PaymentSucceeded(); err != nil {
			return err
		}
		if err := h.repository.SaveIfPaymentCreated(ctx, repairOrder); err != nil {
			return err
		}
		if repairOrder.Timestamps != nil {
			h.metricsPublisher.RecordRepairOrderPhaseDuration("payment_created", time.Since(repairOrder.UpdatedAt).Minutes())
		}
	case StatusFailed, StatusError:
		if err := repairOrder.PaymentFailed(); err != nil {
			return err
		}
		if err := h.repository.SaveIfPaymentCreated(ctx, repairOrder); err != nil {
			return err
		}
	}
	return nil
}
