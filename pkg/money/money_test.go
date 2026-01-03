package money

import (
	"testing"

	"github.com/stretchr/testify/require"
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

func TestAdd(t *testing.T) {
	t.Parallel()

	t.Run("should add two money values", func(t *testing.T) {
		money1, _ := New(1500)
		money2, _ := New(2500)
		result := money1.Add(money2)
		require.Equal(t, Money{Cents: 4000}, result)
	})
}
