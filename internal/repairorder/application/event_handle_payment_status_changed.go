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
	StatusPending    Status = "PENDING"
	StatusProcessing Status = "PROCESSING"
	StatusSucceeded  Status = "SUCCEEDED"
	StatusFailed     Status = "FAILED"
	StatusError      Status = "ERROR"
)

type (
	HandlePaymentStatusChanged struct {
		repository       Repository
		metricsPublisher MetricsPublisher
	}

	paymentTransition struct {
		apply       func() error
		save        func(context.Context, *domain.RepairOrder) error
		metricPhase string
	}
)

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
	transitions := map[Status]paymentTransition{
		StatusPending:    {repairOrder.PaymentCreated, h.repository.SaveIfFinished, "finished"},
		StatusProcessing: {repairOrder.PaymentProcessing, h.repository.SaveIfPaymentCreated, "payment_created"},
		StatusSucceeded:  {repairOrder.PaymentSucceeded, h.repository.SaveIfPaymentProcessing, "payment_processing"},
		StatusFailed:     {repairOrder.PaymentFailed, h.repository.SaveIfPaymentProcessing, ""},
		StatusError:      {repairOrder.PaymentError, h.repository.SaveIfPaymentProcessing, ""},
	}

	transition, ok := transitions[Status(status)]
	if !ok {
		return nil
	}

	if err := transition.apply(); err != nil {
		return err
	}
	if err := transition.save(ctx, repairOrder); err != nil {
		return err
	}
	if transition.metricPhase != "" && repairOrder.Timestamps != nil {
		h.metricsPublisher.RecordRepairOrderPhaseDuration(transition.metricPhase, time.Since(repairOrder.UpdatedAt).Minutes())
	}
	return nil
}
