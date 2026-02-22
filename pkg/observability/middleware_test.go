package observability

import (
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsRepairOrderRoute(t *testing.T) {
	cases := []struct {
		route    string
		expected bool
	}{
		{"/api/repair-orders", true},
		{"/api/repair-orders/:id", true},
		{"/api/customers", false},
		{"/health", false},
		{"", false},
	}

	for _, tc := range cases {
		t.Run(tc.route, func(t *testing.T) {
			assert.Equal(t, tc.expected, isRepairOrderRoute(tc.route))
		})
	}
}

func TestMetricsMiddleware(t *testing.T) {
	t.Run("should pass through when metrics is nil", func(t *testing.T) {
		app := fiber.New()
		app.Use(MetricsMiddleware(nil))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "ok", string(body))
	})
}

func TestRequestLoggingMiddleware(t *testing.T) {
	t.Run("should pass through and log request", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestLoggingMiddleware())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("should set request ID header when missing", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestIDMiddleware())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get(RequestIDHeader))
	})

	t.Run("should preserve existing request ID", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestIDMiddleware())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(RequestIDHeader, "my-custom-id")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, "my-custom-id", resp.Header.Get(RequestIDHeader))
	})
}
