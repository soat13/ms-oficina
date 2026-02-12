package observability

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShutdown(t *testing.T) {
	t.Run("nil components should not panic", func(t *testing.T) {
		assert.NotPanics(t, func() { Shutdown(nil) })
	})

	t.Run("components with nil metrics should not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Shutdown(&Components{Metrics: nil})
		})
	})
}
