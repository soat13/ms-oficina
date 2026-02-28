package observability

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(RequestIDHeader, requestID)
		return c.Next()
	}
}

// ---------------------------------------------------------------------------
// Structured-logging middleware
// ---------------------------------------------------------------------------

func RequestLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		statusCode := c.Response().StatusCode()

		logger := LoggerWithTraceContext(c.UserContext())

		event := logger.Info()
		if statusCode >= 500 {
			event = logger.Error()
		} else if statusCode >= 400 {
			event = logger.Warn()
		}

		event.
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status_code", statusCode).
			Dur("latency", duration).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Msgf("%d - %s %s - %s", statusCode, c.Method(), c.Path(), c.IP())

		return err
	}
}

// ---------------------------------------------------------------------------
// Metrics middleware
// ---------------------------------------------------------------------------

func MetricsMiddleware(metrics *Metrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if metrics == nil {
			return c.Next()
		}

		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		statusCode := c.Response().StatusCode()
		route := c.Route().Path

		metrics.RecordHTTPRequest(c.Method(), route, statusCode, duration)

		return err
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func isRepairOrderRoute(route string) bool {
	return strings.Contains(route, "repair-orders")
}
