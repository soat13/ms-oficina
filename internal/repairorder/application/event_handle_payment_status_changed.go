package application

import (
	"context"
	"time"

	"github.com/soat13/ms-oficina/internal/repairorder/domain"
	"github.com/soat13/ms-oficina/internal/shared/events"
	sharedRepairOrder "github.com/soat13/ms-oficina/internal/shared/repairorder"
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
	repairOrder, err := h.repository.GetById(ctx, evt.RepairOrderID)
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
		StatusPending:    {apply: repairOrder.PaymentCreated, save: h.repository.SaveIfPaymentRequested, metricPhase: "payment_requested"},
		StatusProcessing: {apply: repairOrder.PaymentProcessing, save: h.repository.SaveIfPaymentCreated, metricPhase: "payment_created"},
		StatusFailed:     {apply: repairOrder.PaymentFailed, save: h.repository.SaveIfPaymentProcessing, metricPhase: ""},
		StatusError:      {apply: repairOrder.PaymentError, save: h.repository.SaveIfPaymentProcessingOrFailed, metricPhase: ""},
		StatusSucceeded:  {apply: repairOrder.PaymentSucceeded, save: h.repository.SaveIfPaymentProcessingOrFailed, metricPhase: ""},
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
