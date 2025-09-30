package product

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
// DTOs auxiliares
// -----------------------------------------------------------------------------

type createBody struct {
	Name       string `json:"name"`
	PriceCents int64  `json:"price_cents"`
	Stock      int    `json:"stock,omitempty"`
}

type updateBody struct {
	Name       *string `json:"name,omitempty"`
	PriceCents *int64  `json:"price_cents,omitempty"`
	Stock      *int    `json:"stock,omitempty"`
}

type productJSON struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	PriceCents int64     `json:"price_cents"`
	Stock      int       `json:"stock"`
}

type listResp struct {
	Products []productJSON `json:"products"`
}

type getResp struct {
	Product productJSON `json:"product"`
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func Test_AdminProduct_List_OrderByName(t *testing.T) {
	ensureSetup(t)

	_ = testsupport.ThereIsAProduct(t, env.db, uuid.Nil, "Filtro de Ar", 3000, 5)
	_ = testsupport.ThereIsAProduct(t, env.db, uuid.Nil, "Amortecedor", 20000, 2)

	resp := listProducts(t, 50, 0)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body listResp
	decodeJSON(t, resp, &body)

	idxAmort := indexOfByName(body.Products, "Amortecedor")
	idxFiltro := indexOfByName(body.Products, "Filtro de Ar")

	require.NotEqual(t, -1, idxAmort)
	require.NotEqual(t, -1, idxFiltro)
	require.Less(t, idxAmort, idxFiltro)
}

func indexOfByName(items []productJSON, name string) int {
	for i, s := range items {
		if s.Name == name {
			return i
		}
	}
	return -1
}

func Test_AdminProduct_Update_NotFound(t *testing.T) {
	ensureSetup(t)

	unknown := uuid.New()
	validName := "Produto Válido"

	resp := putUpdateProduct(t, unknown, updateBody{Name: &validName})
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_AdminProduct_Create_OK(t *testing.T) {
	ensureSetup(t)

	payload := createBody{
		Name:       "Filtro de Óleo",
		PriceCents: 4500,
		Stock:      10,
	}
	resp := postCreateProduct(t, payload)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// valida persistência
	var count int
	require.NoError(t,
		env.db.NewRaw(`SELECT COUNT(*) FROM products WHERE name = ? AND price = ? AND stock = ?`,
			payload.Name, payload.PriceCents, payload.Stock).Scan(context.Background(), &count),
	)
	require.Equal(t, 1, count)
}

func Test_AdminProduct_Create_InvalidBody(t *testing.T) {
	ensureSetup(t)

	payload := createBody{Name: "", PriceCents: 0, Stock: -1}
	resp := postCreateProduct(t, payload)
	require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func Test_AdminProduct_Create_DuplicateName(t *testing.T) {
	ensureSetup(t)

	_ = testsupport.ThereIsAProduct(t, env.db, uuid.Nil, "Pneu Aro 17", 200000, 4)

	payload := createBody{Name: "Pneu Aro 17", PriceCents: 210000, Stock: 3}
	resp := postCreateProduct(t, payload)
	require.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func Test_AdminProduct_GetByID_OK(t *testing.T) {
	ensureSetup(t)

	pid := testsupport.ThereIsAProduct(t, env.db, uuid.Nil, "Pastilha de Freio", 8000, 7)
	resp := getProduct(t, pid)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body getResp
	decodeJSON(t, resp, &body)
	require.Equal(t, pid, body.Product.ID)
	require.Equal(t, "Pastilha de Freio", body.Product.Name)
	require.Equal(t, int64(8000), body.Product.PriceCents)
	require.Equal(t, 7, body.Product.Stock)
}

func Test_AdminProduct_GetByID_NotFound(t *testing.T) {
	ensureSetup(t)

	resp := getProduct(t, uuid.New())
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_AdminProduct_Update_OK(t *testing.T) {
	ensureSetup(t)

	pid := testsupport.ThereIsAProduct(t, env.db, uuid.Nil, "Filtro A", 3000, 3)

	newName := "Filtro A Premium"
	newPrice := int64(4500)
	newStock := 8

	resp := putUpdateProduct(t, pid, updateBody{
		Name:       &newName,
		PriceCents: &newPrice,
		Stock:      &newStock,
	})
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	// valida persistência
	var got struct {
		Name  string
		Price int64
		Stock int
	}
	require.NoError(t,
		env.db.NewRaw(`SELECT name, price, stock FROM products WHERE id = ?`, pid).
			Scan(context.Background(), &got),
	)
	require.Equal(t, newName, got.Name)
	require.Equal(t, newPrice, got.Price)
	require.Equal(t, newStock, got.Stock)
}

func Test_AdminProduct_Update_InvalidBody(t *testing.T) {
	ensureSetup(t)

	pid := testsupport.ThereIsAProduct(t, env.db, uuid.Nil, "Filtro Z", 5000, 2)

	neg := int64(-10) // inválido
	resp := putUpdateProduct(t, pid, updateBody{PriceCents: &neg})
	require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func Test_AdminProduct_Delete_OK(t *testing.T) {
	ensureSetup(t)

	pid := testsupport.ThereIsAProduct(t, env.db, uuid.Nil, "Temp P", 1000, 1)
	resp := deleteProduct(t, pid)
	require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	var count int
	require.NoError(t, env.db.NewRaw(`SELECT COUNT(*) FROM products WHERE id = ?`, pid).Scan(context.Background(), &count))
	require.Equal(t, 0, count)
}

// -----------------------------------------------------------------------------
// Helpers HTTP (atenção: rotas no PLURAL)
// -----------------------------------------------------------------------------

func postCreateProduct(t *testing.T, body createBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/admin/products/", bytes.NewReader(bs)) // <-- plural
	req.Header.Set("Content-Type", "application/json")
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getProduct(t *testing.T, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("GET", "/admin/products/"+id.String(), nil) // <-- plural
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func listProducts(t *testing.T, limit, offset int) *http.Response {
	t.Helper()
	url := "/admin/products/" // <-- plural
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
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func putUpdateProduct(t *testing.T, id uuid.UUID, body updateBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/admin/products/"+id.String(), bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func deleteProduct(t *testing.T, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("DELETE", "/admin/products/"+id.String(), nil)
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

// decode helper
func decodeJSON(t *testing.T, resp *http.Response, out any) {
	t.Helper()
	dec := json.NewDecoder(resp.Body)
	require.NoError(t, dec.Decode(out))
}
