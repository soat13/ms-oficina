package observability

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func StartTracer(cfg Config) {
	if !cfg.TraceEnabled {
		log.Info().Msg("Datadog disabled")
		return
	}

	tracer.Start(
		tracer.WithService(cfg.ServiceName),
		tracer.WithEnv(cfg.Environment),
		tracer.WithServiceVersion(cfg.Version),
		tracer.WithRuntimeMetrics(),
		tracer.WithLogStartup(false),
	)

	log.Info().Msg("Datadog configured")
}

func StopTracer() {
	tracer.Stop()
}
