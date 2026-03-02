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
	_ = m.client.Timing("http.request.duration", duration, tags, 1)
	_ = m.client.Incr("http.request.count", tags, 1)
}

// ---------------------------------------------------------------------------
// Repair-order business metrics
// ---------------------------------------------------------------------------

func (m *Metrics) RecordRepairOrderPhaseDuration(phase string, minutes float64) {
	if m == nil || m.client == nil {
		return
	}
	tags := []string{"phase:" + phase}
	_ = m.client.Histogram("repair_order.phase_duration_minutes", minutes, tags, 1)
}
