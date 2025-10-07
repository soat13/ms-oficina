package testsupport

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	authBootstrap "github.com/soat13/fase-1-oficina/internal/bootstrap/auth"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	testauth "github.com/soat13/fase-1-oficina/tests/testsupport/auth"
)

type SetupConfig struct {
	FiberErrorHandler *fiberHelper.ErrorHandler
	Container         *bootstrap.Container
	AuthToken         string
}

func SetupHTTP(t *testing.T, register func(app *fiber.App, c *bootstrap.Container)) *SetupConfig {
	t.Helper()

	testDB := NewTestDB(t)
	fiberApp := fiber.New()
	fiberApp.Use(logger.New())

	container := bootstrap.Build(testDB.DB, fiberApp, nil, nil)

	ensureTestJWTConfig()
	authBootstrap.SetupDefault(container)
	authToken := ensureTestAdminUser(t, fiberApp, container)

	register(fiberApp, container)

	t.Cleanup(func() {
		testDB.Close()
		container.Close()
	})

	return &SetupConfig{
		FiberErrorHandler: container.FiberErrorHandler,
		Container:         container,
		AuthToken:         authToken,
	}
}

func ensureTestAdminUser(t *testing.T, app *fiber.App, container *bootstrap.Container) string {
	adminEmail := "admin@example.com"
	adminPassword := "password123"

	ThereIsAUser(t, container.DB, uuid.Nil, "Test Admin", "71296750043", "11987654999", adminEmail, adminPassword, []string{"manager"})
	authToken := testauth.Authenticate(t, app, adminEmail, adminPassword)
	return authToken
}

func ensureTestJWTConfig() {
	if os.Getenv("JWT_SECRET") == "" {
		_ = os.Setenv("JWT_SECRET", "test-secret")
	}
	if os.Getenv("JWT_EXPIRATION") == "" {
		_ = os.Setenv("JWT_EXPIRATION", "1h")
	}
}
