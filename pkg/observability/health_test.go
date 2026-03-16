package observability

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDBPinger struct {
	err error
}

func (m *mockDBPinger) PingContext(_ context.Context) error {
	return m.err
}

func newTestHealthApp(checker *HealthChecker) *fiber.App {
	app := fiber.New()
	RegisterHealthRoutes(app, checker)
	return app
}

func TestNewHealthChecker(t *testing.T) {
	t.Run("should create health checker", func(t *testing.T) {
		db := &mockDBPinger{}
		hc := NewHealthChecker(db)
		assert.NotNil(t, hc)
	})
}

func TestHealthHandler(t *testing.T) {
	t.Run("should return healthy when db is up", func(t *testing.T) {
		app := newTestHealthApp(NewHealthChecker(&mockDBPinger{}))

		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "healthy")
	})

	t.Run("should return unhealthy when db is down", func(t *testing.T) {
		app := newTestHealthApp(NewHealthChecker(&mockDBPinger{err: errors.New("connection refused")}))

		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "unhealthy")
	})

	t.Run("should return unhealthy when db is nil", func(t *testing.T) {
		app := newTestHealthApp(NewHealthChecker(nil))

		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
	})
}

func TestReadinessHandler(t *testing.T) {
	t.Run("should return ready when db is up", func(t *testing.T) {
		app := newTestHealthApp(NewHealthChecker(&mockDBPinger{}))

		req, _ := http.NewRequest(http.MethodGet, "/health/ready", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "ready")
	})

	t.Run("should return not ready when db is down", func(t *testing.T) {
		app := newTestHealthApp(NewHealthChecker(&mockDBPinger{err: errors.New("timeout")}))

		req, _ := http.NewRequest(http.MethodGet, "/health/ready", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "not_ready")
	})
}

func TestLivenessHandler(t *testing.T) {
	t.Run("should always return alive", func(t *testing.T) {
		app := newTestHealthApp(NewHealthChecker(nil))

		req, _ := http.NewRequest(http.MethodGet, "/health/live", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "alive")
	})
}

func TestStartupHandler(t *testing.T) {
	t.Run("should return started when db is up", func(t *testing.T) {
		app := newTestHealthApp(NewHealthChecker(&mockDBPinger{}))

		req, _ := http.NewRequest(http.MethodGet, "/health/startup", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "started")
	})

	t.Run("should return not_started when db is down", func(t *testing.T) {
		app := newTestHealthApp(NewHealthChecker(&mockDBPinger{err: errors.New("connection refused")}))

		req, _ := http.NewRequest(http.MethodGet, "/health/startup", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "not_started")
	})

	t.Run("should return not_started when db is nil", func(t *testing.T) {
		app := newTestHealthApp(NewHealthChecker(nil))

		req, _ := http.NewRequest(http.MethodGet, "/health/startup", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "not_started")
	})
}

func TestBoolToStatus(t *testing.T) {
	assert.Equal(t, "up", boolToStatus(true))
	assert.Equal(t, "down", boolToStatus(false))
}
