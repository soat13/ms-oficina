package product

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"

	"github.com/soat13/ms-oficina/internal/bootstrap"
	"github.com/soat13/ms-oficina/internal/bootstrap/product"
	"github.com/soat13/ms-oficina/tests/testsupport"
	testauth "github.com/soat13/ms-oficina/tests/testsupport/auth"
)

// -----------------------------------------------------------------------------
// Setup
// -----------------------------------------------------------------------------

func ensureSetup(t *testing.T) *testsupport.SetupConfig {
	t.Helper()

	return testsupport.SetupHTTP(t, func(app *fiber.App, container *bootstrap.Container) {
		product.SetupDefault(container)
	})
}

// -----------------------------------------------------------------------------
// DTOs
// -----------------------------------------------------------------------------

type createBody struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}

type updateBody struct {
	Name  *string `json:"name,omitempty"`
	Price *int64  `json:"price,omitempty"`
	Stock *int    `json:"stock,omitempty"`
}

type productJSON struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price int64     `json:"price"`
	Stock int       `json:"stock"`
}

type listResp struct {
	Products []productJSON `json:"data"`
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func TestCreateProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{
			Name:  "Engine Oil Filter",
			Price: 2500,
			Stock: 50,
		}
		resp := postCreateProduct(t, setup.Container.FiberApp, setup.AuthToken, payload)
		require.Equal(t, fiber.StatusCreated, resp.StatusCode)

		assertProductCreated(t, setup.Container.DB, payload)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{Name: "", Price: 0, Stock: -1}
		resp := postCreateProduct(t, setup.Container.FiberApp, setup.AuthToken, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("DuplicateName", func(t *testing.T) {
		setup := ensureSetup(t)

		_ = testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Brake Pad", 15000, 30)

		payload := createBody{
			Name:  "Brake Pad",
			Price: 16000,
			Stock: 20,
		}
		resp := postCreateProduct(t, setup.Container.FiberApp, setup.AuthToken, payload)
		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})

	t.Run("InvalidPrice", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{
			Name:  "Test Product",
			Price: -100,
			Stock: 10,
		}
		resp := postCreateProduct(t, setup.Container.FiberApp, setup.AuthToken, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("InvalidStock", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{
			Name:  "Test Product",
			Price: 1000,
			Stock: -5,
		}
		resp := postCreateProduct(t, setup.Container.FiberApp, setup.AuthToken, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestGetProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		pid := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Air Filter", 3500, 25)
		resp := getProduct(t, setup.Container.FiberApp, setup.AuthToken, pid)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body productJSON
		decodeJSON(t, resp, &body)
		require.Equal(t, pid, body.ID)
		require.Equal(t, "air filter", body.Name)
		require.Equal(t, int64(3500), body.Price)
		require.Equal(t, 25, body.Stock)
	})

	t.Run("NotFound", func(t *testing.T) {
		setup := ensureSetup(t)

		resp := getProduct(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		pid := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Old Product Name", 5000, 10)

		newName := "New Product Name"
		newPrice := int64(6000)
		newStock := 15

		resp := patchUpdateProduct(t, setup.Container.FiberApp, setup.AuthToken, pid, updateBody{
			Name:  &newName,
			Price: &newPrice,
			Stock: &newStock,
		})
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		var got struct {
			Name  string
			Price int64
			Stock int
		}
		require.NoError(t,
			setup.Container.DB.NewRaw(`SELECT name, price, stock FROM products WHERE id = ?`, pid).
				Scan(context.Background(), &got),
		)
		require.Equal(t, "new product name", got.Name)
		require.Equal(t, newPrice, got.Price)
		require.Equal(t, newStock, got.Stock)
	})

	t.Run("NotFound", func(t *testing.T) {
		setup := ensureSetup(t)

		unknown := uuid.New()
		validName := "Valid Name"

		resp := patchUpdateProduct(t, setup.Container.FiberApp, setup.AuthToken, unknown, updateBody{Name: &validName})
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		setup := ensureSetup(t)

		pid := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Test Product", 1000, 5)

		emptyName := ""
		resp := patchUpdateProduct(t, setup.Container.FiberApp, setup.AuthToken, pid, updateBody{Name: &emptyName})
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("InvalidPrice", func(t *testing.T) {
		setup := ensureSetup(t)

		pid := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Test Product", 1000, 5)

		invalidPrice := int64(-100)
		resp := patchUpdateProduct(t, setup.Container.FiberApp, setup.AuthToken, pid, updateBody{Price: &invalidPrice})
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("InvalidStock", func(t *testing.T) {
		setup := ensureSetup(t)

		pid := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Test Product", 1000, 5)

		invalidStock := -5
		resp := patchUpdateProduct(t, setup.Container.FiberApp, setup.AuthToken, pid, updateBody{Stock: &invalidStock})
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		pid := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Temp Product", 2000, 10)
		resp := deleteProduct(t, setup.Container.FiberApp, setup.AuthToken, pid)
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		var count int
		require.NoError(t,
			setup.Container.DB.NewRaw(`SELECT COUNT(*) FROM products WHERE id = ?`, pid).
				Scan(context.Background(), &count),
		)
		require.Equal(t, 0, count)
	})

	t.Run("NotFound", func(t *testing.T) {
		setup := ensureSetup(t)

		resp := deleteProduct(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})
}

func TestListProducts(t *testing.T) {
	t.Run("OrderByName", func(t *testing.T) {
		setup := ensureSetup(t)

		productBID, productAID := givenProductsOutOfOrder(t, setup.Container.DB)

		// Query with high limit to ensure we get both test products
		response := listProducts(t, setup.Container.FiberApp, setup.AuthToken, 10, 0)
		var body listResp
		decodeJSON(t, response, &body)

		require.Equal(t, fiber.StatusOK, response.StatusCode)

		productBIdx := indexOfByID(body.Products, productBID)
		productAIdx := indexOfByID(body.Products, productAID)

		require.NotEqual(t, -1, productBIdx, "ZZZ Product B should be in the list")
		require.NotEqual(t, -1, productAIdx, "ZZZ Product A should be in the list")
		require.Less(t, productAIdx, productBIdx, "ZZZ Product A should come before ZZZ Product B")
	})
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func assertProductCreated(t *testing.T, db *bun.DB, payload createBody) {
	t.Helper()
	var count int
	require.NoError(t,
		db.NewRaw(
			`SELECT COUNT(*) FROM products WHERE name = ? AND price = ? AND stock = ?`,
			strings.ToLower(payload.Name), payload.Price, payload.Stock,
		).Scan(context.Background(), &count),
	)
	require.Equal(t, 1, count, "Product should be created in database")
}

func givenProductsOutOfOrder(t *testing.T, db *bun.DB) (uuid.UUID, uuid.UUID) {
	t.Helper()
	zProduct := testsupport.ThereIsAProduct(t, db, uuid.Nil, "ZZZ Product B", 1200, 100)
	aProduct := testsupport.ThereIsAProduct(t, db, uuid.Nil, "ZZZ Product A", 3500, 50)
	return zProduct, aProduct
}

func indexOfByID(items []productJSON, ID uuid.UUID) int {
	for i, p := range items {
		if p.ID == ID {
			return i
		}
	}
	return -1
}

func postCreateProduct(t *testing.T, fiberApp *fiber.App, token string, body createBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/admin/products/", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getProduct(t *testing.T, fiberApp *fiber.App, token string, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("GET", "/admin/products/"+id.String(), nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func listProducts(t *testing.T, fiberApp *fiber.App, token string, limit, offset int) *http.Response {
	t.Helper()
	url := "/admin/products/"
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
	testauth.AddAuthHeader(req, token)
	response, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return response
}

func patchUpdateProduct(t *testing.T, fiberApp *fiber.App, token string, id uuid.UUID, body updateBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("PATCH", "/admin/products/"+id.String(), bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func deleteProduct(t *testing.T, fiberApp *fiber.App, token string, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("DELETE", "/admin/products/"+id.String(), nil)
	testauth.AddAuthHeader(req, token)
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
