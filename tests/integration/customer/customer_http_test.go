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

	"github.com/soat13/fase-1-oficina/tests/testsupport"
)

// -----------------------------------------------------------------------------
// DTOs
// -----------------------------------------------------------------------------

type createBody struct {
	Name      string `json:"name"`
	Cellphone string `json:"cellphone"`
	Document  string `json:"document"`
}

type updateBody struct {
	Name      *string `json:"name,omitempty"`
	Cellphone *string `json:"cellphone,omitempty"`
}

type customerJSON struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Cellphone string    `json:"cellphone"`
	Document  string    `json:"document"`
}

type listResp struct {
	Customers []customerJSON `json:"customers"`
}

type getResp struct {
	Customer customerJSON `json:"customer"`
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func TestCreateCustomer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ensureSetup(t)

		payload := createBody{
			Name:      "Ana Silva",
			Cellphone: "11987654321",
			Document:  "11144477735",
		}
		resp := postCreateCustomer(t, payload)
		require.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var count int
		require.NoError(t,
			env.db.NewRaw(
				`SELECT COUNT(*) FROM customers WHERE name = ? AND cellphone = ? AND document = ?`,
				payload.Name, payload.Cellphone, payload.Document,
			).Scan(context.Background(), &count),
		)
		require.Equal(t, 1, count)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		ensureSetup(t)

		payload := createBody{Name: "", Cellphone: "", Document: ""}
		resp := postCreateCustomer(t, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("DuplicateDocument", func(t *testing.T) {
		ensureSetup(t)

		_ = testsupport.ThereIsACustomerWithDocument(t, env.db, uuid.Nil, "João Silva", "11144477735", "CPF", "11987654321")

		payload := createBody{
			Name:      "Maria Silva",
			Cellphone: "11987654322",
			Document:  "11144477735",
		}
		resp := postCreateCustomer(t, payload)
		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})
}

func TestGetCustomer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ensureSetup(t)

		cid := testsupport.ThereIsACustomerWithDocument(t, env.db, uuid.Nil, "Carlos Santos", "98765432100", "CPF", "11987654323")
		resp := getCustomer(t, cid)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body getResp
		decodeJSON(t, resp, &body)
		require.Equal(t, cid, body.Customer.ID)
		require.Equal(t, "Carlos Santos", body.Customer.Name)
		require.Equal(t, "11987654323", body.Customer.Cellphone)
		require.Equal(t, "98765432100", body.Customer.Document)
	})

	t.Run("NotFound", func(t *testing.T) {
		ensureSetup(t)

		resp := getCustomer(t, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestUpdateCustomer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ensureSetup(t)

		cid := testsupport.ThereIsACustomerWithDocument(t, env.db, uuid.Nil, "Pedro Oliveira", "11122233396", "CPF", "11987654324")

		newName := "Pedro Oliveira Santos"
		newCellphone := "11987654325"

		resp := putUpdateCustomer(t, cid, updateBody{
			Name:      &newName,
			Cellphone: &newCellphone,
		})
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var got struct {
			Name      string
			Cellphone string
		}
		require.NoError(t,
			env.db.NewRaw(`SELECT name, cellphone FROM customers WHERE id = ?`, cid).
				Scan(context.Background(), &got),
		)
		require.Equal(t, newName, got.Name)
		require.Equal(t, newCellphone, got.Cellphone)
	})

	t.Run("NotFound", func(t *testing.T) {
		ensureSetup(t)

		unknown := uuid.New()
		validName := "Valid Name"

		resp := putUpdateCustomer(t, unknown, updateBody{Name: &validName})
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		ensureSetup(t)

		cid := testsupport.ThereIsACustomerWithDocument(t, env.db, uuid.Nil, "Test Customer", "52998224725", "CPF", "11987654326")

		emptyName := ""
		resp := putUpdateCustomer(t, cid, updateBody{Name: &emptyName})
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestDeleteCustomer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ensureSetup(t)

		cid := testsupport.ThereIsACustomerWithDocument(t, env.db, uuid.Nil, "Temp Customer", "55566677789", "CPF", "11987654327")
		resp := deleteCustomer(t, cid)
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		var count int
		require.NoError(t,
			env.db.NewRaw(`SELECT COUNT(*) FROM customers WHERE id = ?`, cid).
				Scan(context.Background(), &count),
		)
		require.Equal(t, 0, count)
	})

	t.Run("NotFound", func(t *testing.T) {
		ensureSetup(t)

		resp := deleteCustomer(t, uuid.New())
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})
}

func TestListCustomers(t *testing.T) {
	t.Run("OrderByName", func(t *testing.T) {
		ensureSetup(t)

		mariaID, joaoID := givenCustomersOutOfOrder(t)

		response := listCustomers(t, 2, 0)
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

func givenCustomersOutOfOrder(t *testing.T) (uuid.UUID, uuid.UUID) {
	t.Helper()
	maria := testsupport.ThereIsACustomerWithDocument(t, env.db, uuid.Nil, "Maria Santos", "11144477735", "CPF", "11987654328")
	joao := testsupport.ThereIsACustomerWithDocument(t, env.db, uuid.Nil, "João Silva", "52998224725", "CPF", "11987654329")
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

func postCreateCustomer(t *testing.T, body createBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/admin/customers/", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getCustomer(t *testing.T, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("GET", "/admin/customers/"+id.String(), nil)
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func listCustomers(t *testing.T, limit, offset int) *http.Response {
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
	response, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return response
}

func putUpdateCustomer(t *testing.T, id uuid.UUID, body updateBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/admin/customers/"+id.String(), bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func deleteCustomer(t *testing.T, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("DELETE", "/admin/customers/"+id.String(), nil)
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	require.NoError(t, dec.Decode(v))
}
