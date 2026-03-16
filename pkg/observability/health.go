package observability

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DBPinger interface {
	PingContext(ctx context.Context) error
}

type HealthChecker struct {
	db DBPinger
}

func NewHealthChecker(db DBPinger) *HealthChecker {
	return &HealthChecker{db: db}
}

func RegisterHealthRoutes(app *fiber.App, checker *HealthChecker) {
	app.Get("/health", checker.healthHandler)
	app.Get("/health/ready", checker.readinessHandler)
	app.Get("/health/live", checker.livenessHandler)
	app.Get("/health/startup", checker.startupHandler)
}

func (h *HealthChecker) healthHandler(c *fiber.Ctx) error {
	dbOk := h.checkDB()

	status := "healthy"
	httpStatus := fiber.StatusOK
	if !dbOk {
		status = "unhealthy"
		httpStatus = fiber.StatusServiceUnavailable
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"status":    status,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks": fiber.Map{
			"database": boolToStatus(dbOk),
		},
	})
}

func (h *HealthChecker) readinessHandler(c *fiber.Ctx) error {
	if !h.checkDB() {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "not_ready",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ready",
	})
}

func (h *HealthChecker) livenessHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "alive",
	})
}

func (h *HealthChecker) startupHandler(c *fiber.Ctx) error {
	if !h.checkDB() {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "not_started",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "started",
	})
}

func (h *HealthChecker) checkDB() bool {
	if h.db == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return h.db.PingContext(ctx) == nil
}

func boolToStatus(ok bool) string {
	if ok {
		return "up"
	}
	return "down"
}
