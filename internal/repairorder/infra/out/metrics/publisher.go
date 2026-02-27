package metrics

import "github.com/soat13/fase-1-oficina/pkg/observability"

type MetricsPublisher struct {
	metrics *observability.Metrics
}

func NewPublisher(m *observability.Metrics) *MetricsPublisher {
	return &MetricsPublisher{metrics: m}
}

func (p *MetricsPublisher) IncRepairOrderCreated() {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.IncRepairOrderCreated()
}

func (p *MetricsPublisher) IncRepairOrderStatusChange(fromStatus, toStatus string) {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.IncRepairOrderStatusChange(fromStatus, toStatus)
}

func (p *MetricsPublisher) IncRepairOrderCanceled() {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.IncRepairOrderCanceled()
}

func (p *MetricsPublisher) RecordRepairOrderPhaseDuration(phase string, minutes float64) {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.RecordRepairOrderPhaseDuration(phase, minutes)
}

func (p *MetricsPublisher) RecordRepairOrderExecutionTime(minutes float64) {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.RecordRepairOrderExecutionTime(minutes)
}

func (p *MetricsPublisher) IncIntegrationError(integration, operation string) {
	if p == nil || p.metrics == nil {
		return
	}
	p.metrics.IncIntegrationError(integration, operation)
}
