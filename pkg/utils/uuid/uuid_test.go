package uuid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIDOrNew(t *testing.T) {
	t.Parallel()

	t.Run("should return provided UUID when not nil", func(t *testing.T) {
		providedID := uuid.New()
		result := IDOrNew(providedID)
		require.Equal(t, providedID, result)
	})

	t.Run("should generate new UUID when provided is nil", func(t *testing.T) {
		result := IDOrNew(uuid.Nil)
		require.NotEqual(t, uuid.Nil, result)
		require.NotEqual(t, uuid.UUID{}, result)
	})

	t.Run("should generate different UUIDs on multiple calls", func(t *testing.T) {
		id1 := IDOrNew(uuid.Nil)
		id2 := IDOrNew(uuid.Nil)
		require.NotEqual(t, id1, id2)
	})
}
