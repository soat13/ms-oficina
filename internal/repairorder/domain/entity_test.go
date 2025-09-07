package domain

import (
	"testing"

	"github.com/google/uuid"
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

		assert.Equal(t, StatusReceived, repairOrder.Status)
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

func TestRepairOrder_set(t *testing.T) {
	repairOrder := newRepairOrder()

	t.Run("should return error if not in diagnosis status", func(t *testing.T) {
		err := repairOrder.set(ItemService, uuid.New(), 1)
		assert.ErrorIs(t, err, ErrOperationNotAllowed)
	})

	repairOrder.Status = StatusInDiagnosis

	t.Run("should return error if quantity is negative", func(t *testing.T) {
		err := repairOrder.set(ItemService, uuid.New(), -1)
		assert.ErrorIs(t, err, ErrQuantityMustBeNonNegative)
	})

	t.Run("should return error for invalid item type", func(t *testing.T) {
		err := repairOrder.set("INVALID_TYPE", uuid.New(), 1)
		assert.ErrorIs(t, err, ErrInvalidItemType)
	})

	t.Run("should add service item", func(t *testing.T) {
		serviceID := uuid.New()
		quantity := int64(2)
		currentUpdatedAt := repairOrder.UpdatedAt
		err := repairOrder.set(ItemService, serviceID, quantity)

		require.NoError(t, err)
		assert.Contains(t, repairOrder.Services, serviceID)
		assert.Equal(t, quantity, repairOrder.Services[serviceID])
		assert.True(t, repairOrder.UpdatedAt.After(currentUpdatedAt))

		t.Run("should update service item quantity", func(t *testing.T) {
			newQuantity := int64(5)
			err := repairOrder.set(ItemService, serviceID, newQuantity)

			require.NoError(t, err)
			assert.Equal(t, newQuantity, repairOrder.Services[serviceID])
		})

		t.Run("should add another service item", func(t *testing.T) {
			anotherServiceID := uuid.New()
			anotherQuantity := int64(1)
			err := repairOrder.set(ItemService, anotherServiceID, anotherQuantity)

			require.NoError(t, err)
			assert.Contains(t, repairOrder.Services, anotherServiceID)
			assert.Equal(t, anotherQuantity, repairOrder.Services[anotherServiceID])
			assert.Equal(t, 2, len(repairOrder.Services))
		})

		t.Run("should remove service item when quantity is zero", func(t *testing.T) {
			currentUpdatedAt := repairOrder.UpdatedAt
			err := repairOrder.set(ItemService, serviceID, 0)

			require.NoError(t, err)
			assert.NotContains(t, repairOrder.Services, serviceID)
			assert.True(t, repairOrder.UpdatedAt.After(currentUpdatedAt))
			assert.Equal(t, 1, len(repairOrder.Services))
		})
	})

	t.Run("should add product item", func(t *testing.T) {
		productID := uuid.New()
		quantity := int64(3)
		currentUpdatedAt := repairOrder.UpdatedAt
		err := repairOrder.set(ItemProduct, productID, quantity)

		require.NoError(t, err)
		assert.Contains(t, repairOrder.Products, productID)
		assert.Equal(t, quantity, repairOrder.Products[productID])
		assert.True(t, repairOrder.UpdatedAt.After(currentUpdatedAt))

		t.Run("should update product item quantity", func(t *testing.T) {
			newQuantity := int64(6)
			err := repairOrder.set(ItemProduct, productID, newQuantity)

			require.NoError(t, err)
			assert.Equal(t, newQuantity, repairOrder.Products[productID])
		})

		t.Run("should add another product item", func(t *testing.T) {
			anotherProductID := uuid.New()
			anotherQuantity := int64(4)
			err := repairOrder.set(ItemProduct, anotherProductID, anotherQuantity)

			require.NoError(t, err)
			assert.Contains(t, repairOrder.Products, anotherProductID)
			assert.Equal(t, anotherQuantity, repairOrder.Products[anotherProductID])
			assert.Equal(t, 2, len(repairOrder.Products))
		})

		t.Run("should remove product item when quantity is zero", func(t *testing.T) {
			currentUpdatedAt := repairOrder.UpdatedAt
			err := repairOrder.set(ItemProduct, productID, 0)

			require.NoError(t, err)
			assert.NotContains(t, repairOrder.Products, productID)
			assert.True(t, repairOrder.UpdatedAt.After(currentUpdatedAt))
			assert.Equal(t, 1, len(repairOrder.Products))
		})
	})
}

func TestRepairOrder_SetProductQuantity(t *testing.T) {
	repairOrder := newRepairOrder()
	repairOrder.Status = StatusInDiagnosis

	t.Run("should add product item", func(t *testing.T) {
		productID := uuid.New()
		quantity := int64(3)
		err := repairOrder.SetProductQuantity(productID, quantity)

		require.NoError(t, err)
		assert.Contains(t, repairOrder.Products, productID)
		assert.Equal(t, quantity, repairOrder.Products[productID])
	})
}

func TestRepairOrder_SetServiceQuantity(t *testing.T) {
	repairOrder := newRepairOrder()
	repairOrder.Status = StatusInDiagnosis

	t.Run("should add service item", func(t *testing.T) {
		serviceID := uuid.New()
		quantity := int64(3)
		err := repairOrder.SetServiceQuantity(serviceID, quantity)

		require.NoError(t, err)
		assert.Contains(t, repairOrder.Services, serviceID)
		assert.Equal(t, quantity, repairOrder.Services[serviceID])
	})
}
