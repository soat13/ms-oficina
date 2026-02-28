package observability

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func SetupLogger(cfg Config) {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var writer io.Writer
	if shouldOutputJSON() {
		writer = jsonWriter()
	} else {
		writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}

	level := parseLogLevel(os.Getenv("LOG_LEVEL"))

	log.Logger = zerolog.New(writer).
		With().
		Timestamp().
		Str("dd.service", cfg.ServiceName).
		Str("dd.env", cfg.Environment).
		Str("dd.version", cfg.Version).
		Logger().
		Level(level)
}

func shouldOutputJSON() bool {
	if os.Getenv("DD_LOGS_INJECTION") == "true" {
		return true
	}
	env := os.Getenv("APP_ENV")
	return env != "development" && env != "test" && env != ""
}

func jsonWriter() io.Writer {
	logFile := os.Getenv("DD_LOG_FILE")
	if logFile == "" {
		return os.Stdout
	}

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return os.Stdout
	}
	return io.MultiWriter(os.Stdout, f)
}

func LoggerWithTraceContext(ctx context.Context) zerolog.Logger {
	span, ok := tracer.SpanFromContext(ctx)
	if !ok {
		return log.Logger
	}

	return log.With().
		Str("dd.trace_id", fmt.Sprintf("%d", span.Context().TraceID())).
		Str("dd.span_id", fmt.Sprintf("%d", span.Context().SpanID())).
		Logger()
}

func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
