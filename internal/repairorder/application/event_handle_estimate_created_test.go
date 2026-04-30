package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/repairorder/domain"
	"github.com/soat13/ms-oficina/internal/shared/events"
	"github.com/soat13/ms-oficina/internal/shared/repairorder"
	"github.com/stretchr/testify/assert"
)

func TestHandleEstimateCreated_Execute(t *testing.T) {
	ctx := context.Background()
	roID := uuid.New()

	t.Run("should move repair order to awaiting approval and record metrics", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusDiagnosticsFinished)
		ro.ID = roID

		metrics := &mockMetricsPublisher{}
		repo := &mockRepository{
			getByIdFn: func(_ context.Context, id uuid.UUID) (*domain.RepairOrder, error) {
				return ro, nil
			},
			saveIfDiagnosticsFinishedFn: func(_ context.Context, ro *domain.RepairOrder) error {
				return nil
			},
		}

		handler := NewHandleEstimateCreated(repo, metrics)
		evt := events.EstimateCreated{RepairOrderID: roID}

		err := handler.Execute(ctx, evt)

		assert.NoError(t, err)
		assert.Equal(t, repairorder.StatusAwaitingApproval, ro.Status)
	})

	t.Run("should return error when repository fails to get", func(t *testing.T) {
		metrics := &mockMetricsPublisher{}
		repo := &mockRepository{
			getByIdFn: func(_ context.Context, id uuid.UUID) (*domain.RepairOrder, error) {
				return nil, errors.New("database error")
			},
		}

		handler := NewHandleEstimateCreated(repo, metrics)
		evt := events.EstimateCreated{RepairOrderID: roID}

		err := handler.Execute(ctx, evt)

		assert.Error(t, err)
		assert.Equal(t, 0, metrics.statusChangeCalled)
	})

	t.Run("should return error when status transition is invalid", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusReceived)
		ro.ID = roID

		metrics := &mockMetricsPublisher{}
		repo := &mockRepository{
			getByIdFn: func(_ context.Context, id uuid.UUID) (*domain.RepairOrder, error) {
				return ro, nil
			},
		}

		handler := NewHandleEstimateCreated(repo, metrics)
		evt := events.EstimateCreated{RepairOrderID: roID}

		err := handler.Execute(ctx, evt)

		assert.Error(t, err)
		assert.Equal(t, 0, metrics.statusChangeCalled)
	})

	t.Run("should return error when repository fails to save", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusDiagnosticsFinished)
		ro.ID = roID

		metrics := &mockMetricsPublisher{}
		repo := &mockRepository{
			getByIdFn: func(_ context.Context, id uuid.UUID) (*domain.RepairOrder, error) {
				return ro, nil
			},
			saveIfDiagnosticsFinishedFn: func(_ context.Context, ro *domain.RepairOrder) error {
				return errors.New("save failed")
			},
		}

		handler := NewHandleEstimateCreated(repo, metrics)
		evt := events.EstimateCreated{RepairOrderID: roID}

		err := handler.Execute(ctx, evt)

		assert.Error(t, err)
		assert.Equal(t, 0, metrics.statusChangeCalled)
	})
}
