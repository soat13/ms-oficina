package customer

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
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/customer"
	"github.com/soat13/fase-1-oficina/tests/testsupport"
)

// -----------------------------------------------------------------------------
// Setup
// -----------------------------------------------------------------------------

func ensureSetup(t *testing.T) *testsupport.SetupConfig {
	t.Helper()

	return testsupport.SetupHTTP(t, func(app *fiber.App, container *bootstrap.Container) {
		customer.SetupDefault(container)
	})
}

// -----------------------------------------------------------------------------
// DTOs
// -----------------------------------------------------------------------------

type createBody struct {
	Name        string `json:"name"`
	Document    string `json:"document"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type updateBody struct {
	Name        *string `json:"name,omitempty"`
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

type customerJSON struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Document     string    `json:"document"`
	DocumentType string    `json:"document_type"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number"`
}

type listResp struct {
	Customers []customerJSON `json:"data"`
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func TestCreateCustomer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{
			Name:        "Ana Silva",
			Document:    "11144477735",
			PhoneNumber: "11987654321",
			Email:       "ana@example.com",
		}
		resp := postCreateCustomer(t, setup.FiberApp, payload)
		require.Equal(t, fiber.StatusCreated, resp.StatusCode)

		assertCustomerCreated(t, setup.Container.DB, payload)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{Name: "", PhoneNumber: "", Document: "", Email: ""}
		resp := postCreateCustomer(t, setup.FiberApp, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("DuplicateDocument", func(t *testing.T) {
		setup := ensureSetup(t)

		_ = testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "João Silva", "11144477735", "CPF", "11987654321", "joao@example.com")

		payload := createBody{
			Name:        "Maria Silva",
			PhoneNumber: "11987654322",
			Document:    "11144477735",
			Email:       "maria@example.com",
		}
		resp := postCreateCustomer(t, setup.FiberApp, payload)
		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})

	t.Run("DuplicateDocumentWithFormatting", func(t *testing.T) {
		setup := ensureSetup(t)

		_ = testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "João Silva", "11144477735", "CPF", "11987654321", "joao@example.com")

		payload := createBody{
			Name:        "Maria Silva",
			PhoneNumber: "11987654322",
			Document:    "111.444.777-35",
			Email:       "maria@example.com",
		}
		resp := postCreateCustomer(t, setup.FiberApp, payload)
		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})
}

func TestGetCustomer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		cid := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "Carlos Santos", "98765432100", "CPF", "11987654323", "carlos@example.com")
		resp := getCustomer(t, setup.FiberApp, cid)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body customerJSON
		decodeJSON(t, resp, &body)
		require.Equal(t, cid, body.ID)
		require.Equal(t, "Carlos Santos", body.Name)
		require.Equal(t, "11987654323", body.PhoneNumber)
		require.Equal(t, "98765432100", body.Document)
	})

	t.Run("NotFound", func(t *testing.T) {
		setup := ensureSetup(t)

		resp := getCustomer(t, setup.FiberApp, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestUpdateCustomer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		cid := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "Pedro Oliveira", "85891302071", "CPF", "81989017775", "pedro@example.com")

		newName := "Pedro"
		newPhoneNumber := "81989017776"
		newEmail := "pedro2@example.com"

		resp := putUpdateCustomer(t, setup.FiberApp, cid, updateBody{
			Name:        &newName,
			PhoneNumber: &newPhoneNumber,
			Email:       &newEmail,
		})
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		var got struct {
			Name        string
			Email       string
			PhoneNumber string
		}
		require.NoError(t,
			setup.Container.DB.NewRaw(`SELECT name, email, phone_number FROM customers WHERE id = ?`, cid).
				Scan(context.Background(), &got),
		)
		require.Equal(t, newName, got.Name)
		require.Equal(t, newEmail, got.Email)
		require.Equal(t, newPhoneNumber, got.PhoneNumber)
	})

	t.Run("NotFound", func(t *testing.T) {
		setup := ensureSetup(t)

		unknown := uuid.New()
		validName := "Valid Name"

		resp := putUpdateCustomer(t, setup.FiberApp, unknown, updateBody{Name: &validName})
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		setup := ensureSetup(t)

		cid := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "Test Customer", "52998224725", "CPF", "11987654326", "test@example.com")

		emptyName := ""
		resp := putUpdateCustomer(t, setup.FiberApp, cid, updateBody{Name: &emptyName})
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestDeleteCustomer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		cid := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "Temp Customer", "02383727075", "CPF", "81989017713", "temp@example.com")
		resp := deleteCustomer(t, setup.FiberApp, cid)
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		var count int
		require.NoError(t,
			setup.Container.DB.NewRaw(`SELECT COUNT(*) FROM customers WHERE id = ?`, cid).
				Scan(context.Background(), &count),
		)
		require.Equal(t, 0, count)
	})

	t.Run("NotFound", func(t *testing.T) {
		setup := ensureSetup(t)

		resp := deleteCustomer(t, setup.FiberApp, uuid.New())
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})
}

func TestListCustomers(t *testing.T) {
	t.Run("OrderByName", func(t *testing.T) {
		setup := ensureSetup(t)

		mariaID, joaoID := givenCustomersOutOfOrder(t, setup.Container.DB)

		response := listCustomers(t, setup.FiberApp, 2, 0)
		var body listResp
		decodeJSON(t, response, &body)

		require.Equal(t, fiber.StatusOK, response.StatusCode)

		mariaIdx := indexOfByID(body.Customers, mariaID)
		joaoIdx := indexOfByID(body.Customers, joaoID)

		require.NotEqual(t, -1, mariaIdx, "Maria should be in the list")
		require.NotEqual(t, -1, joaoIdx, "João should be in the list")
		require.Less(t, joaoIdx, mariaIdx, "João should come before Maria")
	})
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func assertCustomerCreated(t *testing.T, db *bun.DB, payload createBody) {
	t.Helper()
	var count int
	require.NoError(t,
		db.NewRaw(
			`SELECT COUNT(*) FROM customers WHERE name = ? AND phone_number = ? AND document = ? AND email = ?`,
			payload.Name, payload.PhoneNumber, payload.Document, payload.Email,
		).Scan(context.Background(), &count),
	)
	require.Equal(t, 1, count, "Customer should be created in database")
}

func givenCustomersOutOfOrder(t *testing.T, db *bun.DB) (uuid.UUID, uuid.UUID) {
	t.Helper()
	maria := testsupport.ThereIsACustomerWithDocument(t, db, uuid.Nil, "Maria Santos", "11144477735", "CPF", "11987654328", "maria@example.com")
	joao := testsupport.ThereIsACustomerWithDocument(t, db, uuid.Nil, "João Silva", "52998224725", "CPF", "11987654329", "joao@example.com")
	return maria, joao
}

func indexOfByID(items []customerJSON, ID uuid.UUID) int {
	for i, c := range items {
		if c.ID == ID {
			return i
		}
	}
	return -1
}

func postCreateCustomer(t *testing.T, fiberApp *fiber.App, body createBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/admin/customers/", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getCustomer(t *testing.T, fiberApp *fiber.App, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("GET", "/admin/customers/"+id.String(), nil)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func listCustomers(t *testing.T, fiberApp *fiber.App, limit, offset int) *http.Response {
	t.Helper()
	url := "/admin/customers/"
	query := ""
	if limit > 0 {
		query += "limit=" + strconv.Itoa(limit)
	}
	if offset > 0 {
		if query != "" {
			query += "&"
		}
		query += "offset=" + strconv.Itoa(offset)
	}
	if query != "" {
		url += "?" + query
	}
	req := httptest.NewRequest("GET", url, nil)
	response, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return response
}

func putUpdateCustomer(t *testing.T, fiberApp *fiber.App, id uuid.UUID, body updateBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("PATCH", "/admin/customers/"+id.String(), bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func deleteCustomer(t *testing.T, fiberApp *fiber.App, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("DELETE", "/admin/customers/"+id.String(), nil)
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
