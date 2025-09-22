package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_MoveToAwaitingApproval_OK(t *testing.T) {
	ro := withStatus(newRO(t), repairorder.StatusInDiagnosis)

	before := ro.UpdatedAt
	require.NoError(t, ro.MoveToAwaitingApproval())

	assert.Equal(t, repairorder.StatusAwaitingApproval, ro.Status)
	assert.True(t, ro.UpdatedAt.After(before))
}

func Test_MoveToAwaitingApproval_Invalid_FromReceived(t *testing.T) {
	ro := withStatus(newRO(t), repairorder.StatusReceived)
	beforeStatus := ro.Status
	beforeUpdated := ro.UpdatedAt

	err := ro.MoveToAwaitingApproval()
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidStatusTransition)

	assert.Equal(t, beforeStatus, ro.Status)
	assert.Equal(t, beforeUpdated, ro.UpdatedAt)
}

func Test_MoveToApproved_OK(t *testing.T) {
	ro := withStatus(newRO(t), repairorder.StatusAwaitingApproval)

	before := ro.UpdatedAt

	require.NoError(t, ro.MoveToApproved())
	assert.Equal(t, repairorder.StatusApproved, ro.Status)
	assert.True(t, ro.UpdatedAt.After(before))
}

func Test_MoveToApproved_Invalid_FromInDiagnosis(t *testing.T) {
	ro := withStatus(newRO(t), repairorder.StatusInDiagnosis)
	beforeStatus := ro.Status
	beforeUpdated := ro.UpdatedAt

	err := ro.MoveToApproved()
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidStatusTransition)

	assert.Equal(t, beforeStatus, ro.Status)
	assert.Equal(t, beforeUpdated, ro.UpdatedAt)
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
