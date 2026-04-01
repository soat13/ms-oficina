package testsupport

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	testauth "github.com/soat13/fase-1-oficina/tests/testsupport/auth"
	fiberHelper "github.com/soat13/oficina-utils/pkg/http/fiber"
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

	container := bootstrap.Build(testDB.DB, fiberApp, nil, nil, nil)

	ensureTestJWTConfig()

	register(fiberApp, container)

	t.Cleanup(func() {
		testDB.Close()
		container.Close()
	})

	return &SetupConfig{
		FiberErrorHandler: container.FiberErrorHandler,
		Container:         container,
		AuthToken: testauth.GenerateToken(),
	}
}

func ensureTestJWTConfig() {
	if os.Getenv("JWT_SECRET") == "" {
		_ = os.Setenv("JWT_SECRET", "test-secret")
	}
	if os.Getenv("JWT_EXPIRATION") == "" {
		_ = os.Setenv("JWT_EXPIRATION", "1h")
	}
}
