package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	"github.com/soat13/fase-1-oficina/internal/shared/repairorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinishExecution_Execute(t *testing.T) {
	ctx := context.Background()
	roID := uuid.New()

	t.Run("should finish execution successfully", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusInExecution)
		metrics := &mockMetricsPublisher{}

		uc := NewFinishExecution(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return ro, nil
			}},
			&mockEventPublisher{},
			metrics,
		)

		err := uc.Execute(ctx, FinishExecutionInput{RepairOrderID: roID})

		require.NoError(t, err)
		assert.Equal(t, repairorder.StatusPaymentRequested, ro.Status)
		assert.NotNil(t, ro.ExecutionTimeMinutes)
		assert.Contains(t, metrics.phaseDurations, "in_execution")
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		repoErr := errors.New("db error")
		uc := NewFinishExecution(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return nil, repoErr
			}},
			&mockEventPublisher{},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, FinishExecutionInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("should return not found when repair order is nil", func(t *testing.T) {
		uc := NewFinishExecution(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return nil, nil
			}},
			&mockEventPublisher{},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, FinishExecutionInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, repairorder.ErrRepairOrderNotFound)
	})

	t.Run("should return error on invalid status transition", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusReceived)
		uc := NewFinishExecution(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return ro, nil
			}},
			&mockEventPublisher{},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, FinishExecutionInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, sharedErrors.ErrInvalidStatusTransaction)
	})

	t.Run("should return error when save fails", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusInExecution)
		saveErr := errors.New("save failed")
		uc := NewFinishExecution(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfInExecutionFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return saveErr
				},
			},
			&mockEventPublisher{},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, FinishExecutionInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, saveErr)
	})
}
