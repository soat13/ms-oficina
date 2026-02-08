package observability

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const RequestIDHeader = "X-Request-ID"

func TracingMiddleware(serviceName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(RequestIDHeader, requestID)

		resourceName := fmt.Sprintf("%s %s", c.Method(), c.Route().Path)

		opts := []tracer.StartSpanOption{
			tracer.ResourceName(resourceName),
			tracer.SpanType(ext.SpanTypeWeb),
			tracer.Tag(ext.HTTPMethod, c.Method()),
			tracer.Tag(ext.HTTPURL, c.OriginalURL()),
			tracer.Tag("http.request_id", requestID),
			tracer.ServiceName(serviceName),
		}

		if spanCtx, err := tracer.Extract(newFiberHeaderCarrier(c)); err == nil {
			opts = append(opts, tracer.ChildOf(spanCtx))
		}

		span, ctx := tracer.StartSpanFromContext(c.UserContext(), "http.request", opts...)
		defer span.Finish()

		c.SetUserContext(ctx)

		err := c.Next()

		statusCode := c.Response().StatusCode()
		span.SetTag(ext.HTTPCode, statusCode)

		if statusCode >= 500 {
			span.SetTag(ext.Error, true)
			span.SetTag("error.message", fmt.Sprintf("HTTP %d", statusCode))
		}
		if err != nil {
			span.SetTag(ext.Error, true)
			span.SetTag("error.message", err.Error())
		}

		return err
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
			Str("request_id", c.Get(RequestIDHeader)).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("route", c.Route().Path).
			Int("status", statusCode).
			Dur("latency", duration).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Msg("request completed")

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

		if statusCode >= 400 && isRepairOrderRoute(route) {
			metrics.IncRepairOrderError(
				c.Method()+" "+route,
				fmt.Sprintf("http_%d", statusCode),
			)
		}

		return err
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func isRepairOrderRoute(route string) bool {
	return strings.Contains(route, "repair-orders")
}

type fiberHeaderCarrier struct {
	ctx *fiber.Ctx
}

func newFiberHeaderCarrier(c *fiber.Ctx) *fiberHeaderCarrier {
	return &fiberHeaderCarrier{ctx: c}
}

func (c *fiberHeaderCarrier) Set(key, val string) {
	c.ctx.Request().Header.Set(key, val)
}

func (c *fiberHeaderCarrier) ForeachKey(handler func(key, val string) error) error {
	var iterErr error
	c.ctx.Request().Header.VisitAll(func(key, value []byte) {
		if iterErr != nil {
			return
		}
		iterErr = handler(string(key), string(value))
	})
	return iterErr
}
