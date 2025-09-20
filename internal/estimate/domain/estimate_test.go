package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepairOrder(t *testing.T) {

	now := time.Now()

	t.Run("should initialize with correct values", func(t *testing.T) {
		estimate, err := NewEstimate(uuid.New(), now, now)
		require.NoError(t, err)
		require.NotNil(t, estimate)

		assert.Equal(t, StatusAwaitingApproval, estimate.Status)
		assert.NotZero(t, estimate.CreatedAt)
		assert.NotZero(t, estimate.UpdatedAt)
		assert.Equal(t, estimate.CreatedAt, estimate.UpdatedAt)
	})

	t.Run("should initialize with invalid repair", func(t *testing.T) {
		estimate, err := NewEstimate(uuid.Nil, now, now)
		assert.Error(t, err)
		assert.Nil(t, estimate)
	})
}
