package vehicle

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
	CustomerID string `json:"customer_id"`
	Plate      string `json:"plate"`
	Model      string `json:"model"`
	Brand      string `json:"brand"`
	Year       int    `json:"year"`
}

type updateBody struct {
	Model *string `json:"model,omitempty"`
	Brand *string `json:"brand,omitempty"`
	Year  *int    `json:"year,omitempty"`
}

type vehicleJSON struct {
	ID         uuid.UUID `json:"id"`
	CustomerID string    `json:"customer_id"`
	Plate      string    `json:"plate"`
	Model      string    `json:"model"`
	Brand      string    `json:"brand"`
	Year       int       `json:"year"`
}

type listResp struct {
	Vehicles []vehicleJSON `json:"vehicles"`
}

type getResp struct {
	Vehicle vehicleJSON `json:"vehicle"`
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func Test_AdminVehicle_Create_OK(t *testing.T) {
	ensureSetup(t)

	customerID := testsupport.ThereIsACustomer(t, env.db, uuid.Nil, "John Doe", "12345678900")

	payload := createBody{
		CustomerID: customerID.String(),
		Plate:      "ABC-1234",
		Model:      "Corolla",
		Brand:      "Toyota",
		Year:       2023,
	}
	resp := postCreateVehicle(t, payload)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var count int
	require.NoError(t,
		env.db.NewRaw(
			`SELECT COUNT(*) FROM vehicles WHERE plate = ?`,
			payload.Plate,
		).Scan(context.Background(), &count),
	)
	require.Equal(t, 1, count)
}

func Test_AdminVehicle_Create_InvalidBody(t *testing.T) {
	ensureSetup(t)

	payload := createBody{Plate: "SHORT"} // Invalid plate
	resp := postCreateVehicle(t, payload)
	require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func Test_AdminVehicle_Create_DuplicatePlate(t *testing.T) {
	ensureSetup(t)
	customerID := testsupport.ThereIsACustomer(t, env.db, uuid.Nil, "Jane Doe", "09876543211")
	_ = testsupport.ThereIsAVehicle(t, env.db, uuid.Nil, customerID, "XYZ-5678", "Honda", "Civic", 2022)

	payload := createBody{Plate: "XYZ-5678"}
	resp := postCreateVehicle(t, payload)
	require.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func Test_AdminVehicle_GetByID_OK(t *testing.T) {
	ensureSetup(t)
	customerID := testsupport.ThereIsACustomer(t, env.db, uuid.Nil, "Get Customer", "11223344556")
	vid := testsupport.ThereIsAVehicle(t, env.db, uuid.Nil, customerID, "GET-0001", "Ford", "Fusion", 2021)
	resp := getVehicle(t, vid)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body getResp
	decodeJSON(t, resp, &body)
	require.Equal(t, vid, body.Vehicle.ID)
	require.Equal(t, "GET-0001", body.Vehicle.Plate)
}

func Test_AdminVehicle_GetByID_NotFound(t *testing.T) {
	ensureSetup(t)

	resp := getVehicle(t, uuid.New())
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_AdminVehicle_Update_OK(t *testing.T) {
	ensureSetup(t)
	customerID := testsupport.ThereIsACustomer(t, env.db, uuid.Nil, "Update Customer", "66778899000")
	vid := testsupport.ThereIsAVehicle(t, env.db, uuid.Nil, customerID, "UPD-0002", "Chevrolet", "Onix", 2020)

	newModel := "Civic"
	newBrand := "Honda"
	newYear := 2024

	resp := putUpdateVehicle(t, vid, updateBody{
		Model: &newModel,
		Brand: &newBrand,
		Year:  &newYear,
	})
	require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	var got struct {
		Model string
		Brand string
		Year  int
	}
	require.NoError(t,
		env.db.NewRaw(`SELECT model, brand, year FROM vehicles WHERE id = ?`, vid).
			Scan(context.Background(), &got),
	)
	require.Equal(t, newModel, got.Model)
	require.Equal(t, newBrand, got.Brand)
	require.Equal(t, newYear, got.Year)
}

func Test_AdminVehicle_Update_NotFound(t *testing.T) {
	ensureSetup(t)
	unknown := uuid.New()
	validModel := "Some Model"
	resp := putUpdateVehicle(t, unknown, updateBody{Model: &validModel})
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_AdminVehicle_Delete_OK(t *testing.T) {
	ensureSetup(t)
	customerID := testsupport.ThereIsACustomer(t, env.db, uuid.Nil, "Delete Customer", "10293847560")
	vid := testsupport.ThereIsAVehicle(t, env.db, uuid.Nil, customerID, "DEL-0003", "VW", "Gol", 2019)
	resp := deleteVehicle(t, vid)
	require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	var count int
	require.NoError(t,
		env.db.NewRaw(`SELECT COUNT(*) FROM vehicles WHERE id = ?`, vid).
			Scan(context.Background(), &count),
	)
	require.Equal(t, 0, count)
}

func Test_AdminVehicle_List_OK(t *testing.T) {
	ensureSetup(t)
	customerID1 := testsupport.ThereIsACustomer(t, env.db, uuid.Nil, "List Customer 1", "11122233344")
	customerID2 := testsupport.ThereIsACustomer(t, env.db, uuid.Nil, "List Customer 2", "55566677788")
	v1 := testsupport.ThereIsAVehicle(t, env.db, uuid.Nil, customerID1, "LST-0004", "Fiat", "Mobi", 2018)
	v2 := testsupport.ThereIsAVehicle(t, env.db, uuid.Nil, customerID2, "LST-0005", "Renault", "Kwid", 2017)

	response := listVehicles(t, 10, 0)
	var body listResp
	decodeJSON(t, response, &body)

	require.Equal(t, fiber.StatusOK, response.StatusCode)
	require.GreaterOrEqual(t, len(body.Vehicles), 2)

	found1 := false
	found2 := false
	for _, v := range body.Vehicles {
		if v.ID == v1 {
			found1 = true
		}
		if v.ID == v2 {
			found2 = true
		}
	}
	require.True(t, found1, "Vehicle 1 should be in the list")
	require.True(t, found2, "Vehicle 2 should be in the list")
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func postCreateVehicle(t *testing.T, body createBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/admin/vehicle/", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getVehicle(t *testing.T, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("GET", "/admin/vehicle/"+id.String(), nil)
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func listVehicles(t *testing.T, limit, offset int) *http.Response {
	t.Helper()
	url := "/admin/vehicle/"
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

func putUpdateVehicle(t *testing.T, id uuid.UUID, body updateBody) *http.Response {
	t.Helper()
	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/admin/vehicle/"+id.String(), bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func deleteVehicle(t *testing.T, id uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("DELETE", "/admin/vehicle/"+id.String(), nil)
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
