package observability

import "os"

type Config struct {
	ServiceName  string
	Environment  string
	Version      string
	AgentHost    string
	StatsDPort   string
	TraceEnabled bool
}

func ConfigFromEnv() Config {
	return Config{
		ServiceName:  getEnvOrDefault("DD_SERVICE", "oficina-api"),
		Environment:  getEnvOrDefault("DD_ENV", "development"),
		Version:      getEnvOrDefault("DD_VERSION", "1.0.0"),
		AgentHost:    getEnvOrDefault("DD_AGENT_HOST", "localhost"),
		StatsDPort:   getEnvOrDefault("DD_DOGSTATSD_PORT", "8125"),
		TraceEnabled: getEnvOrDefault("DD_TRACE_ENABLED", "true") == "true",
	}
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
