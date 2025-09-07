package money

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("should create money with valid cents", func(t *testing.T) {
		money, err := New(1000)
		require.NoError(t, err)
		require.Equal(t, Money{Cents: 1000}, money)
	})

	t.Run("should return error for negative cents", func(t *testing.T) {
		_, err := New(-100)
		require.ErrorIs(t, err, ErrMoneyNegative)
	})

	t.Run("should create money with zero cents", func(t *testing.T) {
		money := Money{}
		require.Equal(t, Money{Cents: 0}, money)
	})
}
