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
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/service"
	"github.com/stretchr/testify/require"

	"github.com/soat13/fase-1-oficina/tests/testsupport"
)

// -----------------------------------------------------------------------------
// Setup
// -----------------------------------------------------------------------------

func ensureSetup(t *testing.T) *testsupport.SetupConfig {
	t.Helper()

	return testsupport.SetupHTTP(t, func(app *fiber.App, container *bootstrap.Container) {
		service.SetupDefault(container)
	})
}

// -----------------------------------------------------------------------------
// DTOs
// -----------------------------------------------------------------------------

type createBody struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type updateBody struct {
	Name  *string `json:"name,omitempty"`
	Price *int64  `json:"price,omitempty"`
}

type serviceJSON struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price int64     `json:"price"`
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

func TestAdminServiceListOrderByName(t *testing.T) {
	setup := ensureSetup(t)

	brakeCheckID, alignID := givenServicesOutOfOrder(t, *setup.Container)

	response := listServices(t, setup.FiberApp, 50, 0)
	var body listResp
	decodeJSON(t, response, &body)

	require.Equal(t, fiber.StatusOK, response.StatusCode)

	brakeIdx := indexOfByID(body.Services, brakeCheckID)
	alignIdx := indexOfByID(body.Services, alignID)

	require.NotEqual(t, -1, brakeIdx, "Brake Check should be in the list")
	require.NotEqual(t, -1, alignIdx, "Alignment should be in the list")
	require.Less(t, alignIdx, brakeIdx, "Alignment should come before Brake Check")
}

func TestAdminServiceUpdateNotFound(t *testing.T) {
	setup := ensureSetup(t)

	unknown := uuid.New()
	validName := "Valid Name"

	resp := putUpdateService(t, setup.FiberApp, unknown, updateBody{Name: &validName})
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestAdminServiceCreateOK(t *testing.T) {
	setup := ensureSetup(t)

	payload := createBody{
		Name:  "Alignment and Balance",
		Price: 12000,
	}
	resp := postCreateService(t, setup.FiberApp, payload)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var count int
	require.NoError(t,
		setup.Container.DB.NewRaw(
			`SELECT COUNT(*) FROM services WHERE name = ? AND price = ?`,
			payload.Name, payload.Price,
		).Scan(context.Background(), &count),
	)
	require.Equal(t, 1, count)
}

func TestAdminServiceCreateInvalidBody(t *testing.T) {
	setup := ensureSetup(t)

	payload := createBody{Name: "", Price: 0}
	resp := postCreateService(t, setup.FiberApp, payload)
	require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func TestAdminServiceACreateDuplicateName(t *testing.T) {
	setup := ensureSetup(t)

	_ = testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Oil Change", 15000)

	payload := createBody{Name: "Oil Change", Price: 20000}
	resp := postCreateService(t, setup.FiberApp, payload)
	require.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestAdminServiceGetByIDOK(t *testing.T) {
	setup := ensureSetup(t)

	sid := testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Rotation", 8000)
	resp := getService(t, setup.FiberApp, sid)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body getResp
	decodeJSON(t, resp, &body)
	require.Equal(t, sid, body.Service.ID)
	require.Equal(t, "Rotation", body.Service.Name)
	require.Equal(t, int64(8000), body.Service.Price)
}

func TestAdminServiceGetByIDNotFound(t *testing.T) {
	setup := ensureSetup(t)

	resp := getService(t, setup.FiberApp, uuid.New())
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestAdminServiceUpdateOK(t *testing.T) {
	setup := ensureSetup(t)

	sid := testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Tire Rotation", 7000)

	newName := "Tire Rotation PRO"
	newPrice := int64(9000)

	resp := putUpdateService(t, setup.FiberApp, sid, updateBody{
		Name:  &newName,
		Price: &newPrice,
	})
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var got struct {
		Name  string
		Price int64
	}
	require.NoError(t,
		setup.Container.DB.NewRaw(`SELECT name, price FROM services WHERE id = ?`, sid).
			Scan(context.Background(), &got),
	)
	require.Equal(t, newName, got.Name)
	require.Equal(t, newPrice, got.Price)
}

func TestAdminServiceUpdateInvalidBody(t *testing.T) {
	setup := ensureSetup(t)

	sid := testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Check A", 5000)

	neg := int64(-10)
	resp := putUpdateService(t, setup.FiberApp, sid, updateBody{Price: &neg})
	require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func TestAdminServiceDeleteOK(t *testing.T) {
	setup := ensureSetup(t)

	sid := testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Temp", 1000)
	resp := deleteService(t, setup.FiberApp, sid)
	require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	var count int
	require.NoError(t,
		setup.Container.DB.NewRaw(`SELECT COUNT(*) FROM services WHERE id = ?`, sid).
			Scan(context.Background(), &count),
	)
	require.Equal(t, 0, count)
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func givenServicesOutOfOrder(t *testing.T, container bootstrap.Container) (uuid.UUID, uuid.UUID) {
	t.Helper()
	brakeCheck := testsupport.ThereIsAService(t, container.DB, uuid.Nil, "Brake Check", 5000)
	align := testsupport.ThereIsAService(t, container.DB, uuid.Nil, "Alignment", 12000)
	return brakeCheck, align
}

func indexOfByID(items []serviceJSON, ID uuid.UUID) int {
	for i, s := range items {
		if s.ID == ID {
			return i
		}
	}
	return -1
}

func postCreateService(t *testing.T, fiberApp *fiber.App, body createBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/admin/services/", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getService(t *testing.T, fiberApp *fiber.App, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("GET", "/admin/services/"+id.String(), nil)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func listServices(t *testing.T, fiberApp *fiber.App, limit, offset int) *http.Response {
	t.Helper()
	url := "/admin/services/"
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

func putUpdateService(t *testing.T, fiberApp *fiber.App, id uuid.UUID, body updateBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/admin/services/"+id.String(), bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func deleteService(t *testing.T, fiberApp *fiber.App, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("DELETE", "/admin/services/"+id.String(), nil)
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
