package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
	"github.com/soat13/fase-1-oficina/internal/shared/repairorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlePaymentStatusChanged_Execute(t *testing.T) {
	ctx := context.Background()
	roID := uuid.New()
	saveFailedError := "save failed"

	newHandler := func(repo Repository, metrics MetricsPublisher) HandlePaymentStatusChanged {
		return NewHandlePaymentStatusChanged(repo, metrics)
	}

	t.Run("returns error when GetById fails", func(t *testing.T) {
		repoErr := errors.New("db error")
		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return nil, repoErr
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "PENDING"})
		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("returns ErrRepairOrderNotFound when GetById returns nil", func(t *testing.T) {
		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return nil, nil
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "PENDING"})
		assert.ErrorIs(t, err, repairorder.ErrRepairOrderNotFound)
	})

	t.Run("StatusPending: sets PaymentURL, transitions Finished->PaymentCreated, records metrics", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusFinished)
		paymentURL := "https://pay.example.com/123"
		metrics := &mockMetricsPublisher{}

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfFinishedFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return nil
				},
			},
			metrics,
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "PENDING", PaymentURL: &paymentURL})

		require.NoError(t, err)
		assert.Equal(t, repairorder.StatusPaymentCreated, ro.Status)
		assert.Equal(t, &paymentURL, ro.PaymentURL)
		assert.Contains(t, metrics.phaseDurations, "finished")
	})

	t.Run("StatusPending: skips UpdatePaymentURL when PaymentURL is nil", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusFinished)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfFinishedFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return nil
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "PENDING", PaymentURL: nil})
		require.NoError(t, err)
		assert.Equal(t, repairorder.StatusPaymentCreated, ro.Status)
		assert.Nil(t, ro.PaymentURL)
	})

	t.Run("StatusPending: returns error on invalid status transition", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusReceived)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "PENDING"})
		assert.ErrorIs(t, err, sharedErrors.ErrInvalidStatusTransaction)
	})

	t.Run("StatusPending: returns error when SaveIfFinished fails", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusFinished)
		saveErr := errors.New(saveFailedError)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfFinishedFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return saveErr
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "PENDING"})
		assert.ErrorIs(t, err, saveErr)
	})

	t.Run("StatusProcessing: sets PaymentURL, transitions PaymentCreated->PaymentProcessing, records metrics", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusPaymentCreated)
		paymentURL := "https://pay.example.com/123"
		metrics := &mockMetricsPublisher{}

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfPaymentCreatedFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return nil
				},
			},
			metrics,
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "PROCESSING", PaymentURL: &paymentURL})

		require.NoError(t, err)
		assert.Equal(t, repairorder.StatusPaymentProcessing, ro.Status)
		assert.Equal(t, &paymentURL, ro.PaymentURL)
		assert.Contains(t, metrics.phaseDurations, "payment_created")
	})

	t.Run("StatusProcessing: skips UpdatePaymentURL when PaymentURL is nil", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusPaymentCreated)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfPaymentCreatedFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return nil
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "PROCESSING", PaymentURL: nil})
		require.NoError(t, err)
		assert.Equal(t, repairorder.StatusPaymentProcessing, ro.Status)
		assert.Nil(t, ro.PaymentURL)
	})

	t.Run("StatusProcessing: returns error on invalid status transition", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusFinished)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "PROCESSING"})
		assert.ErrorIs(t, err, sharedErrors.ErrInvalidStatusTransaction)
	})

	t.Run("StatusProcessing: returns error when SaveIfPaymentCreated fails", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusPaymentCreated)
		saveErr := errors.New(saveFailedError)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfPaymentCreatedFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return saveErr
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "PROCESSING"})
		assert.ErrorIs(t, err, saveErr)
	})

	t.Run("StatusSucceeded: transitions PaymentProcessing->PaymentSucceeded, records metrics", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusPaymentProcessing)
		metrics := &mockMetricsPublisher{}

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfPaymentProcessingFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return nil
				},
			},
			metrics,
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "SUCCEEDED"})

		require.NoError(t, err)
		assert.Equal(t, repairorder.StatusPaymentSucceeded, ro.Status)
		assert.Contains(t, metrics.phaseDurations, "payment_processing")
	})

	t.Run("StatusSucceeded: returns error on invalid status transition", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusReceived)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "SUCCEEDED"})
		assert.ErrorIs(t, err, sharedErrors.ErrInvalidStatusTransaction)
	})

	t.Run("StatusSucceeded: returns error when SaveIfPaymentProcessing fails", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusPaymentProcessing)
		saveErr := errors.New(saveFailedError)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfPaymentProcessingFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return saveErr
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "SUCCEEDED"})
		assert.ErrorIs(t, err, saveErr)
	})

	t.Run("StatusFailed: transitions PaymentProcessing->PaymentFailed", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusPaymentProcessing)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfPaymentProcessingFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return nil
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "FAILED"})

		require.NoError(t, err)
		assert.Equal(t, repairorder.StatusPaymentFailed, ro.Status)
	})

	t.Run("StatusError: transitions PaymentProcessing->PaymentError", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusPaymentProcessing)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfPaymentProcessingFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return nil
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "ERROR"})

		require.NoError(t, err)
		assert.Equal(t, repairorder.StatusPaymentError, ro.Status)
	})

	t.Run("StatusError: returns error on invalid status transition", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusPaymentCreated)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "ERROR"})
		assert.ErrorIs(t, err, sharedErrors.ErrInvalidStatusTransaction)
	})

	t.Run("StatusFailed: returns error when SaveIfPaymentProcessing fails", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusPaymentProcessing)
		saveErr := errors.New(saveFailedError)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfPaymentProcessingFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return saveErr
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "FAILED"})
		assert.ErrorIs(t, err, saveErr)
	})

	t.Run("unknown status: returns nil without modifying order", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusReceived)

		h := newHandler(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
			},
			&mockMetricsPublisher{},
		)

		err := h.Execute(ctx, events.PaymentStatusChanged{ID: roID, Status: "UNKNOWN"})

		assert.NoError(t, err)
		assert.Equal(t, repairorder.StatusReceived, ro.Status)
	})
}
