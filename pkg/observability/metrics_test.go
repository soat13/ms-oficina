package observability

import (
	"testing"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/stretchr/testify/assert"
)

func newTestMetrics() *Metrics {
	client, _ := statsd.New("localhost:0", statsd.WithoutTelemetry())
	return &Metrics{client: client}
}

func TestMetrics_NilSafety(t *testing.T) {
	t.Run("nil receiver should not panic", func(t *testing.T) {
		var m *Metrics
		assert.NotPanics(t, func() { m.RecordHTTPRequest("GET", "/test", 200, time.Millisecond) })
		assert.NotPanics(t, func() { m.IncRepairOrderCreated() })
		assert.NotPanics(t, func() { m.IncRepairOrderStatusChange("a", "b") })
		assert.NotPanics(t, func() { m.RecordRepairOrderExecutionTime(10.0) })
		assert.NotPanics(t, func() { m.IncRepairOrderError("op", "type") })
		assert.NotPanics(t, func() { m.IncRepairOrderCanceled() })
		assert.NotPanics(t, func() { m.IncIntegrationError("pg", "query") })
	})

	t.Run("nil client should not panic", func(t *testing.T) {
		m := &Metrics{client: nil}
		assert.NotPanics(t, func() { m.RecordHTTPRequest("GET", "/test", 200, time.Millisecond) })
		assert.NotPanics(t, func() { m.IncRepairOrderCreated() })
		assert.NotPanics(t, func() { m.IncRepairOrderStatusChange("a", "b") })
		assert.NotPanics(t, func() { m.RecordRepairOrderExecutionTime(10.0) })
		assert.NotPanics(t, func() { m.IncRepairOrderError("op", "type") })
		assert.NotPanics(t, func() { m.IncRepairOrderCanceled() })
		assert.NotPanics(t, func() { m.IncIntegrationError("pg", "query") })
	})
}

func TestMetrics_Close(t *testing.T) {
	t.Run("should close without error", func(t *testing.T) {
		m := newTestMetrics()
		assert.NoError(t, m.Close())
	})

	t.Run("nil receiver should not error", func(t *testing.T) {
		var m *Metrics
		assert.NoError(t, m.Close())
	})

	t.Run("nil client should not error", func(t *testing.T) {
		m := &Metrics{client: nil}
		assert.NoError(t, m.Close())
	})
}

func TestMetrics_RecordHTTPRequest(t *testing.T) {
	m := newTestMetrics()
	assert.NotPanics(t, func() {
		m.RecordHTTPRequest("GET", "/api/repair-orders", 200, 50*time.Millisecond)
	})
}

func TestMetrics_RepairOrderMetrics(t *testing.T) {
	m := newTestMetrics()

	assert.NotPanics(t, func() { m.IncRepairOrderCreated() })
	assert.NotPanics(t, func() { m.IncRepairOrderStatusChange("received", "in_diagnostics") })
	assert.NotPanics(t, func() { m.RecordRepairOrderExecutionTime(30.5) })
	assert.NotPanics(t, func() { m.IncRepairOrderError("create", "validation") })
	assert.NotPanics(t, func() { m.IncRepairOrderCanceled() })
}

func TestMetrics_IntegrationError(t *testing.T) {
	m := newTestMetrics()
	assert.NotPanics(t, func() { m.IncIntegrationError("postgres", "query") })
}
