package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/user"
	"github.com/soat13/fase-1-oficina/tests/testsupport"
)

// -----------------------------------------------------------------------------
// Setup
// -----------------------------------------------------------------------------

func ensureSetup(t *testing.T) *testsupport.SetupConfig {
	t.Helper()

	return testsupport.SetupHTTP(t, func(app *fiber.App, container *bootstrap.Container) {
		user.SetupDefault(container)
	})
}

// -----------------------------------------------------------------------------
// DTOs
// -----------------------------------------------------------------------------

type createBody struct {
	Name        string   `json:"name"`
	Document    string   `json:"document"`
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phone_number"`
	Password    string   `json:"password"`
	Roles       []string `json:"roles"`
}

type updateBody struct {
	Name        *string   `json:"name,omitempty"`
	Email       *string   `json:"email,omitempty"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	Password    *string   `json:"password,omitempty"`
	Roles       *[]string `json:"roles,omitempty"`
}

type userJSON struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Document     string    `json:"document"`
	DocumentType string    `json:"document_type"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number"`
	Roles        []string  `json:"roles"`
}

type listResp struct {
	Users []userJSON `json:"data"`
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func TestCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{
			Name:        "Admin User",
			Document:    "11144477735",
			Email:       "admin@example.com",
			PhoneNumber: "11987654321",
			Password:    "password123",
			Roles:       []string{"attendant"},
		}
		resp := postCreateUser(t, setup.FiberApp, payload)
		require.Equal(t, fiber.StatusCreated, resp.StatusCode)

		assertUserCreated(t, setup.Container.DB, payload)
	})

	t.Run("Success with Multiple Roles", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{
			Name:        "Manager User",
			Document:    "98765432100",
			Email:       "manager@example.com",
			PhoneNumber: "11987654322",
			Password:    "password123",
			Roles:       []string{"manager", "attendant"},
		}
		resp := postCreateUser(t, setup.FiberApp, payload)
		require.Equal(t, fiber.StatusCreated, resp.StatusCode)

		assertUserCreated(t, setup.Container.DB, payload)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{Name: "", Document: "", Email: "", PhoneNumber: "", Password: "", Roles: []string{}}
		resp := postCreateUser(t, setup.FiberApp, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("EmptyRoles", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{
			Name:        "No Roles User",
			Document:    "52998224725",
			Email:       "noroles@example.com",
			PhoneNumber: "11987654323",
			Password:    "password123",
			Roles:       []string{},
		}
		resp := postCreateUser(t, setup.FiberApp, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("InvalidRole", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{
			Name:        "Invalid Role User",
			Document:    "02383727075",
			Email:       "invalidrole@example.com",
			PhoneNumber: "11987654324",
			Password:    "password123",
			Roles:       []string{"invalid_role"},
		}
		resp := postCreateUser(t, setup.FiberApp, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		setup := ensureSetup(t)

		_ = testsupport.ThereIsAUser(t, setup.Container.DB, uuid.Nil, "Existing User", "11144477735", "11987654321", "existing@example.com", "password123", []string{"attendant"})

		payload := createBody{
			Name:        "Another User",
			Document:    "98765432100",
			Email:       "existing@example.com",
			PhoneNumber: "11987654325",
			Password:    "anotherpass",
			Roles:       []string{"manager"},
		}
		resp := postCreateUser(t, setup.FiberApp, payload)
		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})

	t.Run("DuplicateDocument", func(t *testing.T) {
		setup := ensureSetup(t)

		_ = testsupport.ThereIsAUser(t, setup.Container.DB, uuid.Nil, "Existing User", "11144477735", "11987654321", "existing@example.com", "password123", []string{"attendant"})

		payload := createBody{
			Name:        "Another User",
			Document:    "11144477735",
			Email:       "another@example.com",
			PhoneNumber: "11987654326",
			Password:    "anotherpass",
			Roles:       []string{"manager"},
		}
		resp := postCreateUser(t, setup.FiberApp, payload)
		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		uid := testsupport.ThereIsAUser(t, setup.Container.DB, uuid.Nil, "John Mechanic", "85891302071", "11987654321", "user@example.com", "password123", []string{"mechanic"})
		resp := getUser(t, setup.FiberApp, uid)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body userJSON
		err := json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)
		require.Equal(t, uid, body.ID)
		require.Equal(t, "John Mechanic", body.Name)
		require.Equal(t, "user@example.com", body.Email)
		require.Contains(t, body.Roles, "mechanic")
	})

	t.Run("NotFound", func(t *testing.T) {
		setup := ensureSetup(t)

		resp := getUser(t, setup.FiberApp, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		uid := testsupport.ThereIsAUser(t, setup.Container.DB, uuid.Nil, "Old Name", "85891302071", "11987654321", "olduser@example.com", "oldpassword", []string{"mechanic"})

		newName := "New Name"
		newEmail := "newemail@example.com"
		newPhone := "11999999999"
		resp := updateUser(t, setup.FiberApp, uid, updateBody{
			Name:        &newName,
			Email:       &newEmail,
			PhoneNumber: &newPhone,
		})
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		getResp := getUser(t, setup.FiberApp, uid)
		var body userJSON
		err := json.NewDecoder(getResp.Body).Decode(&body)
		require.NoError(t, err)
		require.Equal(t, newName, body.Name)
		require.Equal(t, newEmail, body.Email)
		require.Equal(t, newPhone, body.PhoneNumber)
	})

	t.Run("ChangePassword", func(t *testing.T) {
		setup := ensureSetup(t)

		uid := testsupport.ThereIsAUser(t, setup.Container.DB, uuid.Nil, "Password Test", "52998224725", "11987654321", "pwdtest@example.com", "oldpassword", []string{"attendant"})

		newPassword := "newpassword123"
		resp := updateUser(t, setup.FiberApp, uid, updateBody{
			Password: &newPassword,
		})
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		setup := ensureSetup(t)

		newName := "New Name"
		resp := updateUser(t, setup.FiberApp, uuid.New(), updateBody{
			Name: &newName,
		})
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		setup := ensureSetup(t)

		uid := uuid.New()
		invalidEmail := "not-an-email"
		resp := updateUser(t, setup.FiberApp, uid, updateBody{
			Email: &invalidEmail,
		})
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("UpdateRoles", func(t *testing.T) {
		setup := ensureSetup(t)

		uid := testsupport.ThereIsAUser(t, setup.Container.DB, uuid.Nil, "Role Test", "11144477735", "11987654321", "roletest@example.com", "password123", []string{"attendant"})

		newRoles := []string{"manager", "mechanic"}
		resp := updateUser(t, setup.FiberApp, uid, updateBody{
			Roles: &newRoles,
		})
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		getResp := getUser(t, setup.FiberApp, uid)
		var body userJSON
		err := json.NewDecoder(getResp.Body).Decode(&body)
		require.NoError(t, err)
		require.Equal(t, 2, len(body.Roles))
		require.Contains(t, body.Roles, "manager")
		require.Contains(t, body.Roles, "mechanic")
	})

	t.Run("UpdateRolesEmptyList", func(t *testing.T) {
		setup := ensureSetup(t)

		uid := testsupport.ThereIsAUser(t, setup.Container.DB, uuid.Nil, "Test User", "52998224725", "11987654321", "emptytest@example.com", "password123", []string{"attendant"})

		emptyRoles := []string{}
		resp := updateUser(t, setup.FiberApp, uid, updateBody{
			Roles: &emptyRoles,
		})
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("UpdateRolesInvalid", func(t *testing.T) {
		setup := ensureSetup(t)

		uid := testsupport.ThereIsAUser(t, setup.Container.DB, uuid.Nil, "Invalid Test", "02383727075", "11987654321", "invalidtest@example.com", "password123", []string{"attendant"})

		invalidRoles := []string{"invalid_role"}
		resp := updateUser(t, setup.FiberApp, uid, updateBody{
			Roles: &invalidRoles,
		})
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		uid := testsupport.ThereIsAUser(t, setup.Container.DB, uuid.Nil, "To Delete", "02383727075", "11987654321", "todelete@example.com", "password123", []string{"attendant"})
		resp := deleteUser(t, setup.FiberApp, uid)
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		getResp := getUser(t, setup.FiberApp, uid)
		require.Equal(t, fiber.StatusNotFound, getResp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		setup := ensureSetup(t)

		resp := deleteUser(t, setup.FiberApp, uuid.New())
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})
}

func TestListUsers(t *testing.T) {
	t.Run("OrderByName", func(t *testing.T) {
		setup := ensureSetup(t)

		bob, alice := givenUsers(t, setup.Container.DB)

		resp := listUsers(t, setup.FiberApp, 10, 0)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body listResp
		err := json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(body.Users), 2)

		var foundBob, foundAlice bool
		for _, u := range body.Users {
			if u.ID == bob {
				foundBob = true
				require.Equal(t, "Bob Manager", u.Name)
			}
			if u.ID == alice {
				foundAlice = true
				require.Equal(t, "Alice Mechanic", u.Name)
			}
		}
		require.True(t, foundBob, "Bob not found in list")
		require.True(t, foundAlice, "Alice not found in list")
	})

	t.Run("Pagination", func(t *testing.T) {
		setup := ensureSetup(t)

		givenUsers(t, setup.Container.DB)

		resp := listUsers(t, setup.FiberApp, 1, 0)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body listResp
		err := json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)
		require.Equal(t, 1, len(body.Users))
	})
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func givenUsers(t *testing.T, db *bun.DB) (bob, alice uuid.UUID) {
	t.Helper()
	bob = testsupport.ThereIsAUser(t, db, uuid.Nil, "Bob Manager", "11144477735", "11987654321", "bob@example.com", "password123", []string{"manager"})
	alice = testsupport.ThereIsAUser(t, db, uuid.Nil, "Alice Mechanic", "98765432100", "11987654322", "alice@example.com", "password123", []string{"mechanic"})
	return bob, alice
}

func assertUserCreated(t *testing.T, db *bun.DB, payload createBody) {
	t.Helper()

	var id uuid.UUID
	var name, document, documentType, email, phoneNumber string
	var roles []string

	err := db.NewRaw(`
		SELECT id, name, document, document_type, email, phone_number, roles
		FROM users
		WHERE email = ?
	`, payload.Email).Scan(context.Background(), &id, &name, &document, &documentType, &email, &phoneNumber, pq.Array(&roles))

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)
	require.Equal(t, payload.Name, name)
	require.Equal(t, payload.Email, email)
	require.Equal(t, payload.PhoneNumber, phoneNumber)
	require.Equal(t, len(payload.Roles), len(roles))
}

func postCreateUser(t *testing.T, app *fiber.App, body createBody) *http.Response {
	t.Helper()

	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/admin/users/", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getUser(t *testing.T, app *fiber.App, id uuid.UUID) *http.Response {
	t.Helper()

	req := httptest.NewRequest("GET", "/admin/users/"+id.String(), nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func updateUser(t *testing.T, app *fiber.App, id uuid.UUID, body updateBody) *http.Response {
	t.Helper()

	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest("PATCH", "/admin/users/"+id.String(), bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func deleteUser(t *testing.T, app *fiber.App, id uuid.UUID) *http.Response {
	t.Helper()

	req := httptest.NewRequest("DELETE", "/admin/users/"+id.String(), nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func listUsers(t *testing.T, app *fiber.App, limit, offset int) *http.Response {
	t.Helper()

	req := httptest.NewRequest("GET", "/admin/users/?limit="+strconv.Itoa(limit)+"&offset="+strconv.Itoa(offset), nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	return resp
}
