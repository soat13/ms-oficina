package testsupport

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
)

type SetupConfig struct {
	FiberApp          *fiber.App
	FiberErrorHandler *fiberHelper.ErrorHandler
	Container         *bootstrap.Container
}

func SetupHTTP(t *testing.T, register func(app *fiber.App, c *bootstrap.Container)) *SetupConfig {
	t.Helper()

	testDB := NewTestDB(t)
	fiberApp := fiber.New()
	fiberApp.Use(logger.New())

	container := bootstrap.Build(testDB.DB, fiberApp, fiberHelper.NewErrorHandler(nil))

	register(fiberApp, container)

	t.Cleanup(func() {
		testDB.Close()
		container.Close()
	})

	return &SetupConfig{
		FiberApp:  fiberApp,
		Container: container,
	}
}
