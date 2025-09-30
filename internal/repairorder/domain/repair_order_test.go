package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soat13/fase-1-oficina/internal/shared/errors"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
)

func newRepairOrder() RepairOrder {
	repairOrder, _ := NewRepairOrder(uuid.New(), uuid.New())
	return *repairOrder
}

func TestNewRepairOrder(t *testing.T) {
	t.Run("should initialize with correct values", func(t *testing.T) {
		repairOrder, err := NewRepairOrder(uuid.New(), uuid.New())
		require.NoError(t, err)
		require.NotNil(t, repairOrder)

		assert.Equal(t, repairorder.StatusReceived, repairOrder.Status)
		assert.NotZero(t, repairOrder.CreatedAt)
		assert.NotZero(t, repairOrder.UpdatedAt)
		assert.Equal(t, repairOrder.CreatedAt, repairOrder.UpdatedAt)
	})

	t.Run("should initialize with invalid customer", func(t *testing.T) {
		repairOrder, err := NewRepairOrder(uuid.Nil, uuid.New())
		assert.Error(t, err)
		assert.Nil(t, repairOrder)
	})

	t.Run("should initialize with invalid vehicle", func(t *testing.T) {
		repairOrder, err := NewRepairOrder(uuid.New(), uuid.Nil)
		assert.Error(t, err)
		assert.Nil(t, repairOrder)
	})
}

func TestRepairOrder_touch(t *testing.T) {
	repairOrder := newRepairOrder()
	originalUpdatedAt := repairOrder.UpdatedAt

	repairOrder.Timestamps.Touch()

	assert.True(t, repairOrder.UpdatedAt.After(originalUpdatedAt))
}

func TestRepairOrder_StatusTransitions_Table(t *testing.T) {
	type actionFn func(ro *RepairOrder) error

	moveToAwaiting := func(ro *RepairOrder) error { return ro.MoveToAwaitingApproval() }
	moveToApproved := func(ro *RepairOrder) error { return ro.MoveToApproved() }
	startExecution := func(ro *RepairOrder) error { return ro.StartExecution() }

	cases := []struct {
		name       string
		initial    repairorder.Status
		act        actionFn
		wantStatus repairorder.Status
		wantErr    error
	}{
		// MoveToAwaitingApproval
		{
			name:       "InDiagnosis -> AwaitingApproval (OK)",
			initial:    repairorder.StatusInDiagnostics,
			act:        moveToAwaiting,
			wantStatus: repairorder.StatusAwaitingApproval,
			wantErr:    nil,
		},
		{
			name:       "Received -> AwaitingApproval (ERR)",
			initial:    repairorder.StatusReceived,
			act:        moveToAwaiting,
			wantStatus: repairorder.StatusReceived,
			wantErr:    errors.ErrInvalidStatusTransaction,
		},

		// MoveToApproved
		{
			name:       "AwaitingApproval -> Approved (OK)",
			initial:    repairorder.StatusAwaitingApproval,
			act:        moveToApproved,
			wantStatus: repairorder.StatusApproved,
			wantErr:    nil,
		},
		{
			name:       "InDiagnosis -> Approved (ERR)",
			initial:    repairorder.StatusInDiagnostics,
			act:        moveToApproved,
			wantStatus: repairorder.StatusInDiagnostics,
			wantErr:    errors.ErrInvalidStatusTransaction,
		},

		// StartExecution
		{
			name:       "Approved -> InExecution (OK)",
			initial:    repairorder.StatusApproved,
			act:        startExecution,
			wantStatus: repairorder.StatusInExecution,
			wantErr:    nil,
		},
		{
			name:       "AwaitingApproval -> InExecution (ERR)",
			initial:    repairorder.StatusAwaitingApproval,
			act:        startExecution,
			wantStatus: repairorder.StatusAwaitingApproval,
			wantErr:    errors.ErrInvalidStatusTransaction,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ro := newRO(t)
			withStatus(ro, tc.initial)

			beforeStatus := ro.Status
			beforeUpdated := ro.UpdatedAt

			err := tc.act(ro)

			if tc.wantErr == nil {
				require.NoError(t, err)
				assert.Equal(t, tc.wantStatus, ro.Status)
				assert.True(t, ro.UpdatedAt.After(beforeUpdated))
			} else {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Equal(t, tc.wantStatus, ro.Status)
				assert.Equal(t, beforeUpdated, ro.UpdatedAt)
				assert.Equal(t, beforeStatus, ro.Status)
			}
		})
	}
}

func newRO(t *testing.T) *RepairOrder {
	t.Helper()
	ro, err := NewRepairOrder(uuid.New(), uuid.New())
	require.NoError(t, err)
	require.NotNil(t, ro)
	return ro
}

func withStatus(ro *RepairOrder, st repairorder.Status) *RepairOrder {
	ro.Status = st
	return ro
}
