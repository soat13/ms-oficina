package fiber

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/soat13/ms-oficina/internal/bootstrap"
	"github.com/soat13/ms-oficina/internal/bootstrap/product"
	"github.com/soat13/ms-oficina/tests/testsupport"
	"github.com/stretchr/testify/require"

	testauth "github.com/soat13/ms-oficina/tests/testsupport/auth"
)

func ensureSetup(t *testing.T) *testsupport.SetupConfig {
	t.Helper()

	return testsupport.SetupHTTP(t, func(app *fiber.App, container *bootstrap.Container) {
		product.SetupDefault(container)
	})
}

func TestShouldReturnHttpStatus500(t *testing.T) {

	t.Run("should return 500 for unmapped error", func(t *testing.T) {
		setup := ensureSetup(t)

		setup.Container.FiberApp.Get("/__test__", func(c *fiber.Ctx) error {
			return setup.Container.FiberErrorHandler.Handle(c, errors.New("unmapped error"))
		})

		assertInternalServerError(t, setup.Container.FiberApp, setup.AuthToken)
	})

	t.Run("should return 500 failure in http status code conversion", func(t *testing.T) {
		setup := ensureSetup(t)

		unmappedError := errors.New("error with non-integer status code")
		setup.Container.FiberErrorHandler.ErrorResolver.RegisterError(unmappedError, "non-integer-status-code")
		setup.Container.FiberErrorHandler.ErrorResolver.Resolve(unmappedError)

		setup.Container.FiberApp.Get("/__test__", func(c *fiber.Ctx) error {
			return setup.Container.FiberErrorHandler.Handle(c, unmappedError)
		})

		assertInternalServerError(t, setup.Container.FiberApp, setup.AuthToken)
	})
}

func assertInternalServerError(t *testing.T, fiberApp *fiber.App, authToken string) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/__test__", nil)
	testauth.AddAuthHeader(req, authToken)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
