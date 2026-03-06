package auth

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/soat13/fase-1-oficina/internal/auth/application"
	http2 "github.com/soat13/fase-1-oficina/internal/auth/infra/in/http"
	"github.com/soat13/fase-1-oficina/internal/auth/infra/out/db"
	"github.com/soat13/fase-1-oficina/internal/auth/infra/out/jwt"
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
)

type Config struct {
	Secret            string
	TokenTTL          time.Duration
	ProtectedPrefixes []string
	PublicRoutes      []http2.PublicRoute
}

func SetupDefault(container *bootstrap.Container) {
	cfg := Config{
		Secret:            os.Getenv("JWT_SECRET"),
		TokenTTL:          parseDuration(os.Getenv("JWT_EXPIRATION")),
		ProtectedPrefixes: []string{"/"},
		PublicRoutes: []http2.PublicRoute{
			{
				Method: fiber.MethodPost,
				Path:   "/auth/login",
			},
			{
				Method: fiber.MethodGet,
				Path:   "/docs",
			},
			{
				Method: fiber.MethodGet,
				Path:   "/openapi.yaml",
			},
			{
				Method: fiber.MethodGet,
				Path:   "/favicon.ico",
			},
		},
	}
	Setup(container, cfg)
}

func Setup(container *bootstrap.Container, cfg Config) {
	if cfg.Secret == "" {
		log.Fatal("JWT_SECRET must be configured")
	}
	if cfg.TokenTTL <= 0 {
		cfg.TokenTTL = time.Hour
	}
	if len(cfg.ProtectedPrefixes) == 0 {
		cfg.ProtectedPrefixes = []string{"/"}
	}
	if len(cfg.PublicRoutes) == 0 {
		cfg.PublicRoutes = []http2.PublicRoute{
			{
				Method: fiber.MethodPost,
				Path:   "/auth/login",
			},
		}
	}

	tokenService, err := jwt.NewTokenService(cfg.Secret, cfg.TokenTTL)
	if err != nil {
		log.Fatalf("failed to build token service: %v", err)
	}

	userReader := db.NewBunUserReader(container.DB)
	authenticate := application.NewAuthenticateUser(userReader, tokenService)
	handler := http2.NewHandler(authenticate, container.Validator, container.FiberErrorHandler)

	middleware := http2.NewMiddleware(tokenService, container.FiberErrorHandler, cfg.ProtectedPrefixes, cfg.PublicRoutes)
	container.FiberApp.Use(middleware.Handle)

	if os.Getenv("APP_ENV") != "production" {
		http2.Register(container.FiberApp, handler)
	}
}

func parseDuration(value string) time.Duration {
	if value == "" {
		return 0
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("⚠️  invalid JWT_EXPIRATION value %q, falling back to default\n", value)
		return 0
	}
	return d
}
