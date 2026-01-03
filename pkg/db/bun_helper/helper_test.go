package bun_helper

import (
	"errors"
	"testing"
)

func TestShouldHandleDeleteError(t *testing.T) {
	t.Parallel()

	t.Run("foreign key violation", func(t *testing.T) {
		err := HandleDeleteError(
			errors.New("pq: insert or update on table \"orders\" violates foreign key constraint \"orders_user_id_fkey\" (SQLSTATE 23503)"),
		)

		if !errors.Is(err, ErrResourceInUse) {
			t.Fatalf("expected ErrResourceInUse, got %v", err)
		}
	})

	t.Run("unmapped error", func(t *testing.T) {
		err := HandleDeleteError(nil)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})
}
