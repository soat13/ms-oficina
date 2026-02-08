package observability

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Components struct {
	Metrics       *Metrics
	HealthChecker *HealthChecker
}

func Setup(app *fiber.App, db DBPinger) *Components {
	SetupLogger()

	ddCfg := ConfigFromEnv()

	StartTracer(ddCfg)

	metrics, err := NewMetrics(ddCfg)
	if err != nil {
		log.Warn().Err(err).Msg("Metrics client unavailable")
	}

	app.Use(TracingMiddleware(ddCfg.ServiceName))
	app.Use(RequestLoggingMiddleware())
	app.Use(MetricsMiddleware(metrics))

	healthChecker := NewHealthChecker(db)
	RegisterHealthRoutes(app, healthChecker)

	log.Info().Msg("Observability configured")

	return &Components{
		Metrics:       metrics,
		HealthChecker: healthChecker,
	}
}

func Shutdown(components *Components) {
	if components == nil {
		return
	}

	StopTracer()
	if components.Metrics != nil {
		_ = components.Metrics.Close()
	}

	log.Info().Msg("Observability shutdown")
}
