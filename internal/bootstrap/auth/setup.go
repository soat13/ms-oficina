package auth

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	authApp "github.com/soat13/fase-1-oficina/internal/auth/application"
	authHTTP "github.com/soat13/fase-1-oficina/internal/auth/infra/http"
	authJWT "github.com/soat13/fase-1-oficina/internal/auth/infra/jwt"
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	userDB "github.com/soat13/fase-1-oficina/internal/user/infra/db"
)

type Config struct {
	Secret            string
	TokenTTL          time.Duration
	ProtectedPrefixes []string
	PublicRoutes      []authHTTP.PublicRoute
}

func SetupDefault(container *bootstrap.Container) {
	cfg := Config{
		Secret:            os.Getenv("JWT_SECRET"),
		TokenTTL:          parseDuration(os.Getenv("JWT_EXPIRATION")),
		ProtectedPrefixes: []string{"/"},
		PublicRoutes: []authHTTP.PublicRoute{
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
		cfg.PublicRoutes = []authHTTP.PublicRoute{
			{
				Method: fiber.MethodPost,
				Path:   "/auth/login",
			},
		}
	}

	tokenService, err := authJWT.NewTokenService(cfg.Secret, cfg.TokenTTL)
	if err != nil {
		log.Fatalf("failed to build token service: %v", err)
	}

	userRepo := userDB.NewBunUserRepository(container.DB)
	authenticate := authApp.NewAuthenticateUser(userRepo, tokenService)
	handler := authHTTP.NewHandler(authenticate, container.Validator, container.FiberErrorHandler)

	middleware := authHTTP.NewMiddleware(tokenService, container.FiberErrorHandler, cfg.ProtectedPrefixes, cfg.PublicRoutes)
	container.FiberApp.Use(middleware.Handle)

	authHTTP.Register(container.FiberApp, handler)
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
