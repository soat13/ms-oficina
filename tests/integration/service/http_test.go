package service

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

	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/service"
	"github.com/soat13/fase-1-oficina/tests/testsupport"
	testauth "github.com/soat13/fase-1-oficina/tests/testsupport/auth"
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
	Data []serviceJSON `json:"data"`
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func TestCreateService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{
			Name:  "Alignment and Balance",
			Price: 12000,
		}
		resp := postCreateService(t, setup.Container.FiberApp, setup.AuthToken, payload)
		require.Equal(t, fiber.StatusCreated, resp.StatusCode)

		assertServiceCreated(t, setup.Container.DB, payload)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		setup := ensureSetup(t)

		payload := createBody{Name: "", Price: 0}
		resp := postCreateService(t, setup.Container.FiberApp, setup.AuthToken, payload)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("DuplicateName", func(t *testing.T) {
		setup := ensureSetup(t)

		_ = testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Oil Change", 15000)

		payload := createBody{Name: "Oil Change", Price: 20000}
		resp := postCreateService(t, setup.Container.FiberApp, setup.AuthToken, payload)
		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})
}

func TestGetService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		sid := testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Rotation", 8000)
		resp := getService(t, setup.Container.FiberApp, setup.AuthToken, sid)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body serviceJSON
		decodeJSON(t, resp, &body)
		require.Equal(t, sid, body.ID)
		require.Equal(t, "rotation", body.Name)
		require.Equal(t, int64(8000), body.Price)
	})

	t.Run("NotFound", func(t *testing.T) {
		setup := ensureSetup(t)

		resp := getService(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestUpdateService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		sid := testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Tire Rotation", 7000)

		expectedName := "Tire Rotation PRO"
		expectedPrice := int64(9000)

		resp := patchUpdateService(t, setup.Container.FiberApp, setup.AuthToken, sid, updateBody{
			Name:  &expectedName,
			Price: &expectedPrice,
		})

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		var got struct {
			Name  string
			Price int64
		}
		require.NoError(t,
			setup.Container.DB.NewRaw(`SELECT name, price FROM services WHERE id = ?`, sid).
				Scan(context.Background(), &got),
		)
		require.Equal(t, "tire rotation pro", got.Name)
		require.Equal(t, expectedPrice, got.Price)
	})

	t.Run("NotFound", func(t *testing.T) {
		setup := ensureSetup(t)

		unknown := uuid.New()
		validName := "Valid Name"

		resp := patchUpdateService(t, setup.Container.FiberApp, setup.AuthToken, unknown, updateBody{Name: &validName})
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		setup := ensureSetup(t)

		sid := testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Check A", 5000)

		neg := int64(-10)
		resp := patchUpdateService(t, setup.Container.FiberApp, setup.AuthToken, sid, updateBody{Price: &neg})
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestDeleteService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup := ensureSetup(t)

		sid := testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Temp", 1000)
		resp := deleteService(t, setup.Container.FiberApp, setup.AuthToken, sid)
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		var count int
		require.NoError(t,
			setup.Container.DB.NewRaw(`SELECT COUNT(*) FROM services WHERE id = ?`, sid).
				Scan(context.Background(), &count),
		)
		require.Equal(t, 0, count)
	})
}

func TestListServices(t *testing.T) {
	t.Run("OrderByName", func(t *testing.T) {
		setup := ensureSetup(t)

		brakeCheckID, alignID := givenServicesOutOfOrder(t, *setup.Container)

		response := listServices(t, setup.Container.FiberApp, setup.AuthToken, 50, 0)
		var body listResp
		decodeJSON(t, response, &body)

		require.Equal(t, fiber.StatusOK, response.StatusCode)

		brakeIdx := indexOfByID(body.Data, brakeCheckID)
		alignIdx := indexOfByID(body.Data, alignID)

		require.NotEqual(t, -1, brakeIdx)
		require.NotEqual(t, -1, alignIdx)
		require.Less(t, alignIdx, brakeIdx)
	})
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func assertServiceCreated(t *testing.T, db *bun.DB, payload createBody) {
	t.Helper()
	var count int
	require.NoError(t,
		db.NewRaw(
			`SELECT COUNT(*) FROM services WHERE name = ? AND price = ?`,
			strings.ToLower(payload.Name), payload.Price,
		).Scan(context.Background(), &count),
	)
	require.Equal(t, 1, count)
}

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

func postCreateService(t *testing.T, fiberApp *fiber.App, token string, body createBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/admin/services/", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getService(t *testing.T, fiberApp *fiber.App, token string, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("GET", "/admin/services/"+id.String(), nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func listServices(t *testing.T, fiberApp *fiber.App, token string, limit, offset int) *http.Response {
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
	testauth.AddAuthHeader(req, token)
	response, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return response
}

func patchUpdateService(t *testing.T, fiberApp *fiber.App, token string, id uuid.UUID, body updateBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("PATCH", "/admin/services/"+id.String(), bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func deleteService(t *testing.T, fiberApp *fiber.App, token string, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("DELETE", "/admin/services/"+id.String(), nil)
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
