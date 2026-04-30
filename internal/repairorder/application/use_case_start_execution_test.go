package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/repairorder/domain"
	sharedErrors "github.com/soat13/ms-oficina/internal/shared/errors"
	"github.com/soat13/ms-oficina/internal/shared/repairorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartExecution_Execute(t *testing.T) {
	ctx := context.Background()
	roID := uuid.New()

	t.Run("should start execution successfully", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusApproved)
		metrics := &mockMetricsPublisher{}

		uc := NewStartExecution(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return ro, nil
			}},
			metrics,
		)

		err := uc.Execute(ctx, StartExecutionInput{RepairOrderID: roID})

		require.NoError(t, err)
		assert.Equal(t, repairorder.StatusInExecution, ro.Status)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		repoErr := errors.New("db error")
		uc := NewStartExecution(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return nil, repoErr
			}},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, StartExecutionInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("should return not found when repair order is nil", func(t *testing.T) {
		uc := NewStartExecution(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return nil, nil
			}},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, StartExecutionInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, repairorder.ErrRepairOrderNotFound)
	})

	t.Run("should return error on invalid status transition", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusReceived)
		uc := NewStartExecution(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return ro, nil
			}},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, StartExecutionInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, sharedErrors.ErrInvalidStatusTransaction)
	})

	t.Run("should return error when save fails", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusApproved)
		saveErr := errors.New("save failed")
		uc := NewStartExecution(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfApprovedFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return saveErr
				},
			},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, StartExecutionInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, saveErr)
	})
}
