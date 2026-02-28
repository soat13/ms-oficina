package observability

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	fibertrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gofiber/fiber.v2"
)

type Components struct {
	Metrics       *Metrics
	HealthChecker *HealthChecker
}

func Setup(app *fiber.App, db DBPinger) *Components {
	ddCfg := ConfigFromEnv()

	SetupLogger(ddCfg)

	StartTracer(ddCfg)

	metrics, err := NewMetrics(ddCfg)
	if err != nil {
		log.Warn().Err(err).Msg("Metrics client unavailable")
	}

	app.Use(RequestIDMiddleware())
	app.Use(fibertrace.Middleware(fibertrace.WithServiceName(ddCfg.ServiceName)))
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
