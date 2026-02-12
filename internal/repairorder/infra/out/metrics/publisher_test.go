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
		assert.NotPanics(t, func() { p.IncRepairOrderCreated() })
		assert.NotPanics(t, func() { p.IncRepairOrderStatusChange("a", "b") })
		assert.NotPanics(t, func() { p.IncRepairOrderCanceled() })
		assert.NotPanics(t, func() { p.RecordRepairOrderExecutionTime(10.5) })
		assert.NotPanics(t, func() { p.IncIntegrationError("test", "create") })
	})

	t.Run("nil metrics field should not panic", func(t *testing.T) {
		p := &MetricsPublisher{metrics: nil}
		assert.NotPanics(t, func() { p.IncRepairOrderCreated() })
		assert.NotPanics(t, func() { p.IncRepairOrderStatusChange("a", "b") })
		assert.NotPanics(t, func() { p.IncRepairOrderCanceled() })
		assert.NotPanics(t, func() { p.RecordRepairOrderExecutionTime(10.5) })
		assert.NotPanics(t, func() { p.IncIntegrationError("test", "create") })
	})
}

func TestMetricsPublisher_IncRepairOrderCreated(t *testing.T) {
	t.Run("should not panic with valid metrics", func(t *testing.T) {
		cfg := observability.Config{
			ServiceName: "test", Environment: "test", Version: "1.0.0",
			AgentHost: "localhost", StatsDPort: "8125", TraceEnabled: false,
		}
		metrics, _ := observability.NewMetrics(cfg)
		publisher := NewPublisher(metrics)
		
		assert.NotPanics(t, func() { publisher.IncRepairOrderCreated() })
	})
}

func TestMetricsPublisher_IncRepairOrderStatusChange(t *testing.T) {
	t.Run("should not panic with valid metrics and params", func(t *testing.T) {
		cfg := observability.Config{
			ServiceName: "test", Environment: "test", Version: "1.0.0",
			AgentHost: "localhost", StatsDPort: "8125", TraceEnabled: false,
		}
		metrics, _ := observability.NewMetrics(cfg)
		publisher := NewPublisher(metrics)
		
		assert.NotPanics(t, func() { 
			publisher.IncRepairOrderStatusChange("received", "in_diagnostics") 
		})
	})
}

func TestMetricsPublisher_IncRepairOrderCanceled(t *testing.T) {
	t.Run("should not panic with valid metrics", func(t *testing.T) {
		cfg := observability.Config{
			ServiceName: "test", Environment: "test", Version: "1.0.0",
			AgentHost: "localhost", StatsDPort: "8125", TraceEnabled: false,
		}
		metrics, _ := observability.NewMetrics(cfg)
		publisher := NewPublisher(metrics)
		
		assert.NotPanics(t, func() { publisher.IncRepairOrderCanceled() })
	})
}

func TestMetricsPublisher_RecordRepairOrderExecutionTime(t *testing.T) {
	t.Run("should not panic with valid metrics and time", func(t *testing.T) {
		cfg := observability.Config{
			ServiceName: "test", Environment: "test", Version: "1.0.0",
			AgentHost: "localhost", StatsDPort: "8125", TraceEnabled: false,
		}
		metrics, _ := observability.NewMetrics(cfg)
		publisher := NewPublisher(metrics)
		
		assert.NotPanics(t, func() { 
			publisher.RecordRepairOrderExecutionTime(45.5) 
		})
	})

	t.Run("should not panic with zero time", func(t *testing.T) {
		cfg := observability.Config{
			ServiceName: "test", Environment: "test", Version: "1.0.0",
			AgentHost: "localhost", StatsDPort: "8125", TraceEnabled: false,
		}
		metrics, _ := observability.NewMetrics(cfg)
		publisher := NewPublisher(metrics)
		
		assert.NotPanics(t, func() { 
			publisher.RecordRepairOrderExecutionTime(0) 
		})
	})
}

func TestMetricsPublisher_IncIntegrationError(t *testing.T) {
	t.Run("should not panic with valid metrics and params", func(t *testing.T) {
		cfg := observability.Config{
			ServiceName: "test", Environment: "test", Version: "1.0.0",
			AgentHost: "localhost", StatsDPort: "8125", TraceEnabled: false,
		}
		metrics, _ := observability.NewMetrics(cfg)
		publisher := NewPublisher(metrics)
		
		assert.NotPanics(t, func() { 
			publisher.IncIntegrationError("payment-service", "process_payment") 
		})
	})
}
