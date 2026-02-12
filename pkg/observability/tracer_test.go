package observability

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartTracer(t *testing.T) {
	t.Run("should not start when trace disabled", func(t *testing.T) {
		cfg := Config{
			ServiceName:  "test",
			Environment:  "test",
			Version:      "0.0.1",
			TraceEnabled: false,
		}
		assert.NotPanics(t, func() { StartTracer(cfg) })
	})

	t.Run("should start and stop without error", func(t *testing.T) {
		cfg := Config{
			ServiceName:  "test",
			Environment:  "test",
			Version:      "0.0.1",
			TraceEnabled: true,
		}
		assert.NotPanics(t, func() {
			StartTracer(cfg)
			StopTracer()
		})
	})
}

func TestStopTracer(t *testing.T) {
	t.Run("should not panic when called without start", func(t *testing.T) {
		assert.NotPanics(t, func() { StopTracer() })
	})
}
