package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/repairorder/domain"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	"github.com/soat13/fase-1-oficina/internal/shared/repairorder"
	"github.com/soat13/oficina-utils/pkg/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRepairOrderWithStatus(status repairorder.Status) *domain.RepairOrder {
	now := time.Now()
	s := status
	ro, _ := domain.NewRepairOrder(uuid.New(), uuid.New(), uuid.New(), &s, &now, &now)
	return ro
}

func TestStartDiagnostics_Execute(t *testing.T) {
	ctx := context.Background()
	roID := uuid.New()

	t.Run("should start diagnostics successfully", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusReceived)
		metrics := &mockMetricsPublisher{}

		uc := NewStartDiagnostics(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return ro, nil
			}},
			metrics,
		)

		err := uc.Execute(ctx, StartDiagnosticsInput{RepairOrderID: roID})

		require.NoError(t, err)
		assert.Equal(t, repairorder.StatusInDiagnostics, ro.Status)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		repoErr := errors.New("db error")
		uc := NewStartDiagnostics(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return nil, repoErr
			}},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, StartDiagnosticsInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("should return not found when repair order is nil", func(t *testing.T) {
		uc := NewStartDiagnostics(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return nil, nil
			}},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, StartDiagnosticsInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, repairorder.ErrRepairOrderNotFound)
	})

	t.Run("should return error on invalid status transition", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusFinished)
		uc := NewStartDiagnostics(
			&mockRepository{getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
				return ro, nil
			}},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, StartDiagnosticsInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, sharedErrors.ErrInvalidStatusTransaction)
	})

	t.Run("should return error when save fails", func(t *testing.T) {
		ro := newRepairOrderWithStatus(repairorder.StatusReceived)
		saveErr := errors.New("save failed")
		uc := NewStartDiagnostics(
			&mockRepository{
				getByIdFn: func(_ context.Context, _ uuid.UUID) (*domain.RepairOrder, error) {
					return ro, nil
				},
				saveIfReceivedFn: func(_ context.Context, _ *domain.RepairOrder) error {
					return saveErr
				},
			},
			&mockMetricsPublisher{},
		)

		err := uc.Execute(ctx, StartDiagnosticsInput{RepairOrderID: roID})

		assert.ErrorIs(t, err, saveErr)
	})
}

func TestNewRepairOrderWithStatus(t *testing.T) {
	ro := newRepairOrderWithStatus(repairorder.StatusApproved)
	require.NotNil(t, ro)
	assert.Equal(t, repairorder.StatusApproved, ro.Status)
	assert.IsType(t, &entity.Timestamps{}, ro.Timestamps)
}
