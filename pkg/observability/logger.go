package observability

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func SetupLogger() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	env := os.Getenv("APP_ENV")

	var writer io.Writer
	if env == "development" || env == "test" || env == "" {
		writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	} else {
		writer = os.Stdout
	}

	level := parseLogLevel(os.Getenv("LOG_LEVEL"))

	log.Logger = zerolog.New(writer).
		With().
		Timestamp().
		Logger().
		Level(level)
}

func LoggerWithTraceContext(ctx context.Context) zerolog.Logger {
	span, ok := tracer.SpanFromContext(ctx)
	if !ok {
		return log.Logger
	}

	return log.With().
		Uint64("dd.trace_id", span.Context().TraceID()).
		Uint64("dd.span_id", span.Context().SpanID()).
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
