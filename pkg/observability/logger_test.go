package observability

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestParseLogLevel(t *testing.T) {
	cases := []struct {
		input    string
		expected zerolog.Level
	}{
		{"trace", zerolog.TraceLevel},
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"", zerolog.InfoLevel},
		{"unknown", zerolog.InfoLevel},
	}

	for _, tc := range cases {
		t.Run("level_"+tc.input, func(t *testing.T) {
			assert.Equal(t, tc.expected, parseLogLevel(tc.input))
		})
	}
}

func TestSetupLogger(t *testing.T) {
	cfg := Config{
		ServiceName: "test-svc",
		Environment: "test",
		Version:     "0.0.1",
	}

	t.Run("should not panic in development", func(t *testing.T) {
		t.Setenv("APP_ENV", "development")
		t.Setenv("LOG_LEVEL", "debug")
		t.Setenv("DD_LOGS_INJECTION", "")
		assert.NotPanics(t, func() { SetupLogger(cfg) })
	})

	t.Run("should not panic in production", func(t *testing.T) {
		t.Setenv("APP_ENV", "production")
		t.Setenv("LOG_LEVEL", "info")
		t.Setenv("DD_LOGS_INJECTION", "")
		assert.NotPanics(t, func() { SetupLogger(cfg) })
	})

	t.Run("should not panic with empty env", func(t *testing.T) {
		t.Setenv("APP_ENV", "")
		t.Setenv("LOG_LEVEL", "")
		t.Setenv("DD_LOGS_INJECTION", "")
		assert.NotPanics(t, func() { SetupLogger(cfg) })
	})

	t.Run("should not panic with DD_LOGS_INJECTION enabled", func(t *testing.T) {
		t.Setenv("APP_ENV", "development")
		t.Setenv("LOG_LEVEL", "info")
		t.Setenv("DD_LOGS_INJECTION", "true")
		assert.NotPanics(t, func() { SetupLogger(cfg) })
	})
}

func TestLoggerWithTraceContext(t *testing.T) {
	t.Run("should return default logger when no span in context", func(t *testing.T) {
		ctx := t.Context()
		logger := LoggerWithTraceContext(ctx)
		assert.NotNil(t, logger)
	})
}
