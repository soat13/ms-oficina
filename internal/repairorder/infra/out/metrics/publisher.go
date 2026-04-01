package metrics

import "github.com/soat13/oficina-utils/pkg/observability"

type MetricsPublisher struct {
	metrics *observability.Metrics
}

func NewPublisher(m *observability.Metrics) *MetricsPublisher {
	return &MetricsPublisher{metrics: m}
}

func (p *MetricsPublisher) RecordRepairOrderPhaseDuration(phase string, minutes float64) {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.RecordRepairOrderPhaseDuration(phase, minutes)
}
