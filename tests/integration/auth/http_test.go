package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/auth"
	"github.com/soat13/fase-1-oficina/tests/testsupport"
)

// -----------------------------------------------------------------------------
// Setup
// -----------------------------------------------------------------------------

func ensureSetup(t *testing.T) *testsupport.SetupConfig {
	t.Helper()

	return testsupport.SetupHTTP(t, func(app *fiber.App, container *bootstrap.Container) {
		auth.SetupDefault(container)
	})
}

// -----------------------------------------------------------------------------
// DTOs
// -----------------------------------------------------------------------------

type authenticateBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authenticateResponse struct {
	AccessToken string          `json:"access_token"`
	TokenType   string          `json:"token_type"`
	ExpiresIn   int64           `json:"expires_in"`
	User        userInformation `json:"user"`
}

type userInformation struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Roles []string  `json:"roles"`
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func TestAuthenticate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		userID := testsupport.ThereIsAUser(
			t,
			setup.Container.DB,
			uuid.Nil,
			"John Mechanic",
			"85891302071",
			"11987654321",
			"mechanic@example.com",
			"password123",
			[]string{"mechanic"},
		)

		payload := authenticateBody{
			Email:    "mechanic@example.com",
			Password: "password123",
		}

		resp := postAuthenticate(t, setup.Container.FiberApp, payload)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body authenticateResponse
		decodeJSON(t, resp, &body)

		require.NotEmpty(t, body.AccessToken)
		require.Equal(t, "Bearer", body.TokenType)
		require.Equal(t, userID, body.User.ID)
		require.Equal(t, "John Mechanic", body.User.Name)
		require.Equal(t, "mechanic@example.com", body.User.Email)
		require.Contains(t, body.User.Roles, "mechanic")
	})

	t.Run("Success with Multiple Roles", func(t *testing.T) {
		setup := ensureSetup(t)

		userID := testsupport.ThereIsAUser(
			t,
			setup.Container.DB,
			uuid.Nil,
			"Manager User",
			"52998224725",
			"11987654322",
			"manager@example.com",
			"password456",
			[]string{"manager", "attendant"},
		)

		payload := authenticateBody{
			Email:    "manager@example.com",
			Password: "password456",
		}

		resp := postAuthenticate(t, setup.Container.FiberApp, payload)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body authenticateResponse
		decodeJSON(t, resp, &body)

		require.NotEmpty(t, body.AccessToken)
		require.Equal(t, "Bearer", body.TokenType)
		require.Equal(t, userID, body.User.ID)
		require.Contains(t, body.User.Roles, "manager")
		require.Contains(t, body.User.Roles, "attendant")
	})

	t.Run("Invalid Credentials - Wrong Password", func(t *testing.T) {
		setup := ensureSetup(t)

		testsupport.ThereIsAUser(
			t,
			setup.Container.DB,
			uuid.Nil,
			"Test User",
			"02383727075",
			"11987654323",
			"test@example.com",
			"correctpassword",
			[]string{"attendant"},
		)

		payload := authenticateBody{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		resp := postAuthenticate(t, setup.Container.FiberApp, payload)
		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Invalid Credentials - User Not Found", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := authenticateBody{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		resp := postAuthenticate(t, setup.Container.FiberApp, payload)
		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Invalid Body - Empty Email", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := authenticateBody{
			Email:    "",
			Password: "password123",
		}

		resp := postAuthenticate(t, setup.Container.FiberApp, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Invalid Body - Empty Password", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := authenticateBody{
			Email:    "test@example.com",
			Password: "",
		}

		resp := postAuthenticate(t, setup.Container.FiberApp, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Invalid Body - Invalid Email Format", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := authenticateBody{
			Email:    "not-an-email",
			Password: "password123",
		}

		resp := postAuthenticate(t, setup.Container.FiberApp, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		setup := ensureSetup(t)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := setup.Container.FiberApp.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func postAuthenticate(t *testing.T, fiberApp *fiber.App, body authenticateBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	require.NoError(t, dec.Decode(v))
}
