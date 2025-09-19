package service

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
	Currency   string `json:"currency,omitempty"`
}

type updateBody struct {
	Name       *string `json:"name,omitempty"`
	PriceCents *int64  `json:"price_cents,omitempty"`
	Currency   *string `json:"currency,omitempty"`
}

type serviceJSON struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	PriceCents int64     `json:"price_cents"`
	Currency   string    `json:"currency"`
}

type listResp struct {
	Services []serviceJSON `json:"services"`
}

type getResp struct {
	Service serviceJSON `json:"service"`
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func Test_AdminService_List_OrderByName(t *testing.T) {
	ensureSetup(t)

	_ = testsupport.ThereIsAService(t, env.db, uuid.Nil, "Brake Check", 5000)
	_ = testsupport.ThereIsAService(t, env.db, uuid.Nil, "Alignment", 12000)

	resp := listServices(t, 50, 0)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body listResp
	decodeJSON(t, resp, &body)

	// A lista contém seeds + criados no teste. Confira a ordem relativa dos criados.
	idxAlignment := indexOfByName(body.Services, "Alignment")
	idxBrake := indexOfByName(body.Services, "Brake Check")

	require.NotEqual(t, -1, idxAlignment, "Alignment deve estar na lista")
	require.NotEqual(t, -1, idxBrake, "Brake Check deve estar na lista")
	require.Less(t, idxAlignment, idxBrake, "Alignment deve vir antes de Brake Check")
}

func indexOfByName(items []serviceJSON, name string) int {
	for i, s := range items {
		if s.Name == name {
			return i
		}
	}
	return -1
}

func Test_AdminService_Update_NotFound(t *testing.T) {
	ensureSetup(t)

	unknown := uuid.New()
	validName := "Valid Name" // precisa ser válido para não disparar 422 de validação

	resp := putUpdateService(t, unknown, updateBody{Name: &validName})
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_AdminService_Create_OK(t *testing.T) {
	ensureSetup(t)

	payload := createBody{
		Name:       "Alignment",
		PriceCents: 12000,
		// Currency vazio => default BRL no domínio
	}
	resp := postCreateService(t, payload)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// valida persistência
	var count int
	require.NoError(t,
		env.db.NewRaw(`SELECT COUNT(*) FROM services WHERE name = ? AND price = ? AND currency = 'BRL'`,
			payload.Name, payload.PriceCents).Scan(context.Background(), &count),
	)
	require.Equal(t, 1, count)
}

func Test_AdminService_Create_InvalidBody(t *testing.T) {
	ensureSetup(t)

	// name vazio e price <= 0
	payload := createBody{Name: "", PriceCents: 0}
	resp := postCreateService(t, payload)
	require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func Test_AdminService_Create_DuplicateName(t *testing.T) {
	ensureSetup(t)

	// seed direto no DB usando helper
	_ = testsupport.ThereIsAService(t, env.db, uuid.Nil, "Oil Change", 15000)

	payload := createBody{Name: "Oil Change", PriceCents: 20000}
	resp := postCreateService(t, payload)
	require.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func Test_AdminService_GetByID_OK(t *testing.T) {
	ensureSetup(t)

	sid := testsupport.ThereIsAService(t, env.db, uuid.Nil, "Rotation", 8000)
	resp := getService(t, sid)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body getResp
	decodeJSON(t, resp, &body)
	require.Equal(t, sid, body.Service.ID)
	require.Equal(t, "Rotation", body.Service.Name)
	require.Equal(t, int64(8000), body.Service.PriceCents)
	require.Equal(t, "BRL", body.Service.Currency)
}

func Test_AdminService_GetByID_NotFound(t *testing.T) {
	ensureSetup(t)

	resp := getService(t, uuid.New())
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_AdminService_Update_OK(t *testing.T) {
	ensureSetup(t)

	sid := testsupport.ThereIsAService(t, env.db, uuid.Nil, "Tire Rotation", 7000)

	newName := "Tire Rotation PRO"
	newPrice := int64(9000)
	newCurrency := "USD"

	resp := putUpdateService(t, sid, updateBody{
		Name:       &newName,
		PriceCents: &newPrice,
		Currency:   &newCurrency,
	})
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	// valida persistência
	var got struct {
		Name     string
		Price    int64
		Currency string
	}
	require.NoError(t,
		env.db.NewRaw(`SELECT name, price, currency FROM services WHERE id = ?`, sid).
			Scan(context.Background(), &got),
	)
	require.Equal(t, newName, got.Name)
	require.Equal(t, newPrice, got.Price)
	require.Equal(t, "USD", got.Currency) // domínio uppercasa
}

func Test_AdminService_Update_InvalidBody(t *testing.T) {
	ensureSetup(t)

	sid := testsupport.ThereIsAService(t, env.db, uuid.Nil, "Check A", 5000)

	neg := int64(-10) // inválido
	resp := putUpdateService(t, sid, updateBody{PriceCents: &neg})
	require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func Test_AdminService_Delete_OK(t *testing.T) {
	ensureSetup(t)

	sid := testsupport.ThereIsAService(t, env.db, uuid.Nil, "Temp", 1000)
	resp := deleteService(t, sid)
	require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	// verifica exclusão
	var count int
	require.NoError(t, env.db.NewRaw(`SELECT COUNT(*) FROM services WHERE id = ?`, sid).Scan(context.Background(), &count))
	require.Equal(t, 0, count)
}

// -----------------------------------------------------------------------------
// Helpers HTTP (atenção: rotas no PLURAL)
// -----------------------------------------------------------------------------

func postCreateService(t *testing.T, body createBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/admin/services/", bytes.NewReader(bs)) // <-- plural
	req.Header.Set("Content-Type", "application/json")
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getService(t *testing.T, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("GET", "/admin/services/"+id.String(), nil) // <-- plural
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func listServices(t *testing.T, limit, offset int) *http.Response {
	t.Helper()
	url := "/admin/services/" // <-- plural
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

func putUpdateService(t *testing.T, id uuid.UUID, body updateBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/admin/services/"+id.String(), bytes.NewReader(bs)) // <-- plural
	req.Header.Set("Content-Type", "application/json")
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func deleteService(t *testing.T, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("DELETE", "/admin/services/"+id.String(), nil) // <-- plural
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

// -----------------------------------------------------------------------------
// Utils
// -----------------------------------------------------------------------------

func decodeJSON(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	require.NoError(t, dec.Decode(v))
}

func strPtr(s string) *string { return &s }
