package observability

import (
	"fmt"
	"os"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

type Metrics struct {
	client statsd.ClientInterface
}

func NewMetrics(cfg Config) (*Metrics, error) {
	addr := os.Getenv("DD_DOGSTATSD_URL")
	if addr == "" {
		addr = fmt.Sprintf("%s:%s", cfg.AgentHost, cfg.StatsDPort)
	}

	client, err := statsd.New(addr,
		statsd.WithNamespace("oficina."),
		statsd.WithTags([]string{
			"service:" + cfg.ServiceName,
			"env:" + cfg.Environment,
			"version:" + cfg.Version,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("statsd client: %w", err)
	}

	return &Metrics{client: client}, nil
}

func (m *Metrics) Close() error {
	if m == nil || m.client == nil {
		return nil
	}
	return m.client.Close()
}

// ---------------------------------------------------------------------------
// HTTP metrics
// ---------------------------------------------------------------------------

func (m *Metrics) RecordHTTPRequest(method, route string, statusCode int, duration time.Duration) {
	if m == nil || m.client == nil {
		return
	}
	tags := []string{
		"method:" + method,
		"route:" + route,
		fmt.Sprintf("status_code:%d", statusCode),
		fmt.Sprintf("status_class:%dxx", statusCode/100),
	}

	ms := float64(duration.Microseconds()) / 1000.0
	_ = m.client.Histogram("http.request.duration", ms, tags, 1)
	_ = m.client.Incr("http.request.count", tags, 1)
}

// ---------------------------------------------------------------------------
// Repair-order business metrics
// ---------------------------------------------------------------------------

func (m *Metrics) IncRepairOrderCreated() {
	if m == nil || m.client == nil {
		return
	}
	_ = m.client.Incr("repair_order.created", nil, 1)
}

func (m *Metrics) IncRepairOrderStatusChange(fromStatus, toStatus string) {
	if m == nil || m.client == nil {
		return
	}
	tags := []string{
		"from_status:" + fromStatus,
		"to_status:" + toStatus,
	}
	_ = m.client.Incr("repair_order.status_change", tags, 1)
}

func (m *Metrics) RecordRepairOrderExecutionTime(minutes float64) {
	if m == nil || m.client == nil {
		return
	}
	_ = m.client.Histogram("repair_order.execution_time_minutes", minutes, nil, 1)
}

func (m *Metrics) IncRepairOrderError(operation, errType string) {
	if m == nil || m.client == nil {
		return
	}
	tags := []string{
		"operation:" + operation,
		"error_type:" + errType,
	}
	_ = m.client.Incr("repair_order.error", tags, 1)
}

func (m *Metrics) IncRepairOrderCanceled() {
	if m == nil || m.client == nil {
		return
	}
	_ = m.client.Incr("repair_order.canceled", nil, 1)
}

// ---------------------------------------------------------------------------
// Integration error metrics
// ---------------------------------------------------------------------------

func (m *Metrics) IncIntegrationError(integration, operation string) {
	if m == nil || m.client == nil {
		return
	}
	tags := []string{
		"integration:" + integration,
		"operation:" + operation,
	}
	_ = m.client.Incr("integration.error", tags, 1)
}
