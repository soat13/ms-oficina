package metrics

import (
	"testing"

	"github.com/soat13/fase-1-oficina/pkg/observability"
	"github.com/stretchr/testify/assert"
)

func TestNewPublisher(t *testing.T) {
	t.Run("should create publisher with nil metrics", func(t *testing.T) {
		publisher := NewPublisher(nil)
		assert.NotNil(t, publisher)
		assert.Nil(t, publisher.metrics)
	})

	t.Run("should create publisher with metrics", func(t *testing.T) {
		cfg := observability.Config{
			ServiceName:  "test-service",
			Environment:  "test",
			Version:      "1.0.0",
			AgentHost:    "localhost",
			StatsDPort:   "8125",
			TraceEnabled: false,
		}
		metrics, _ := observability.NewMetrics(cfg)
		publisher := NewPublisher(metrics)
		assert.NotNil(t, publisher)
		assert.NotNil(t, publisher.metrics)
	})
}

func TestMetricsPublisher_NilSafety(t *testing.T) {
	t.Run("nil receiver should not panic", func(t *testing.T) {
		var p *MetricsPublisher
		assert.NotPanics(t, func() { p.RecordRepairOrderPhaseDuration("in_diagnostics", 5.0) })
	})

	t.Run("nil metrics field should not panic", func(t *testing.T) {
		p := &MetricsPublisher{metrics: nil}
		assert.NotPanics(t, func() { p.RecordRepairOrderPhaseDuration("in_execution", 15.0) })
	})
}
