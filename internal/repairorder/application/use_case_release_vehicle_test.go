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

func TestReleaseVehicle_Execute(t *testing.T) {
	ctx := context.Background()
	roID := uuid.New()

	t.Run("should release vehicle successfully", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusPaymentSucceeded)
		metrics := &mockMetricsPublisher{}

		uc := NewReleaseVehicle(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return ro, nil
			}},
			metrics,
		)

		err := uc.Execute(ctx, ReleaseVehicleInput{RepairOrderID: roID})

		require.NoError(t, err)
		assert.Equal(t, repairorder.StatusReleased, ro.Status)
		assert.Contains(t, metrics.phaseDurations, "payment_succeeded")
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		repoErr := errors.New("db error")
		uc := NewReleaseVehicle(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return nil, repoErr
			}},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, ReleaseVehicleInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("should return not found when repair order is nil", func(t *testing.T) {
		uc := NewReleaseVehicle(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return nil, nil
			}},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, ReleaseVehicleInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, repairorder.ErrRepairOrderNotFound)
	})

	t.Run("should return error on invalid status transition", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusReceived)
		uc := NewReleaseVehicle(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return ro, nil
			}},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, ReleaseVehicleInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, sharedErrors.ErrInvalidStatusTransaction)
	})

	t.Run("should return error when save fails", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusPaymentSucceeded)
		saveErr := errors.New("save failed")
		uc := NewReleaseVehicle(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfPaymentSucceededFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return saveErr
				},
			},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, ReleaseVehicleInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, saveErr)
	})
}
