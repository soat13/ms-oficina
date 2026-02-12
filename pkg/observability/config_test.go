package observability

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFromEnv(t *testing.T) {
	t.Run("should return defaults when no env vars set", func(t *testing.T) {
		t.Setenv("DD_SERVICE", "")
		t.Setenv("DD_ENV", "")
		t.Setenv("DD_VERSION", "")
		t.Setenv("DD_AGENT_HOST", "")
		t.Setenv("DD_DOGSTATSD_PORT", "")
		t.Setenv("DD_TRACE_ENABLED", "")

		cfg := ConfigFromEnv()

		assert.Equal(t, "oficina-api", cfg.ServiceName)
		assert.Equal(t, "development", cfg.Environment)
		assert.Equal(t, "1.0.0", cfg.Version)
		assert.Equal(t, "localhost", cfg.AgentHost)
		assert.Equal(t, "8125", cfg.StatsDPort)
		assert.True(t, cfg.TraceEnabled)
	})

	t.Run("should read from env vars", func(t *testing.T) {
		t.Setenv("DD_SERVICE", "my-svc")
		t.Setenv("DD_ENV", "production")
		t.Setenv("DD_VERSION", "2.0.0")
		t.Setenv("DD_AGENT_HOST", "agent.local")
		t.Setenv("DD_DOGSTATSD_PORT", "9999")
		t.Setenv("DD_TRACE_ENABLED", "false")

		cfg := ConfigFromEnv()

		assert.Equal(t, "my-svc", cfg.ServiceName)
		assert.Equal(t, "production", cfg.Environment)
		assert.Equal(t, "2.0.0", cfg.Version)
		assert.Equal(t, "agent.local", cfg.AgentHost)
		assert.Equal(t, "9999", cfg.StatsDPort)
		assert.False(t, cfg.TraceEnabled)
	})
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Run("should return env value when set", func(t *testing.T) {
		t.Setenv("TEST_KEY", "value")
		assert.Equal(t, "value", getEnvOrDefault("TEST_KEY", "fallback"))
	})

	t.Run("should return fallback when env is empty", func(t *testing.T) {
		t.Setenv("TEST_KEY", "")
		assert.Equal(t, "fallback", getEnvOrDefault("TEST_KEY", "fallback"))
	})

	t.Run("should return fallback when env is unset", func(t *testing.T) {
		assert.Equal(t, "fallback", getEnvOrDefault("NONEXISTENT_KEY_12345", "fallback"))
	})
}
