package vehicle

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/vehicle"
	testauth "github.com/soat13/fase-1-oficina/tests/testsupport/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soat13/fase-1-oficina/tests/testsupport"
)

func ensureSetup(t *testing.T) *testsupport.SetupConfig {
	t.Helper()

	return testsupport.SetupHTTP(t, func(app *fiber.App, container *bootstrap.Container) {
		vehicle.SetupDefault(container)
	})
}

func TestVehicleCreate(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		customerID := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "John Doe", "11144477735", "CPF", "11987654321", "john@example.com")

		body := createBody{
			CustomerID: customerID,
			Plate:      "ABC1234",
			Brand:      "Toyota",
			Model:      "Corolla",
			Year:       2020,
		}

		resp := postCreateVehicle(t, setup, body)
		require.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid Body", func(t *testing.T) {
		body := createBody{
			CustomerID: uuid.Nil,
			Plate:      "",
			Brand:      "",
			Model:      "",
			Year:       0,
		}

		resp := postCreateVehicle(t, setup, body)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Duplicate Plate", func(t *testing.T) {
		customerID := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "Jane Doe", "52998224725", "CPF", "11987654322", "jane@example.com")
		plate := "XYZ5678"

		body1 := createBody{
			CustomerID: customerID,
			Plate:      plate,
			Brand:      "Honda",
			Model:      "Civic",
			Year:       2019,
		}
		resp1 := postCreateVehicle(t, setup, body1)
		require.Equal(t, fiber.StatusCreated, resp1.StatusCode)

		body2 := createBody{
			CustomerID: customerID,
			Plate:      plate,
			Brand:      "Ford",
			Model:      "Focus",
			Year:       2021,
		}
		resp2 := postCreateVehicle(t, setup, body2)
		require.Equal(t, fiber.StatusConflict, resp2.StatusCode)
	})
}

func TestVehicleUpdate(t *testing.T) {
	setup := ensureSetup(t)
	t.Run("Success", func(t *testing.T) {
		customerID := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "John Doe", "11144477735", "CPF", "11987654321", "john@example.com")
		vehicleID := testsupport.ThereIsAVehicle(t, setup.Container.DB, uuid.Nil, customerID, "ABC1234", "Toyota", "Corolla", 2020)

		body := updateBody{
			Brand: stringPtr("Honda"),
			Model: stringPtr("Civic"),
			Year:  intPtr(2021),
		}

		resp := patchUpdateVehicle(t, setup, vehicleID, body)
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("Not Found", func(t *testing.T) {
		body := updateBody{
			Brand: stringPtr("Honda"),
		}

		resp := patchUpdateVehicle(t, setup, uuid.New(), body)
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestVehicleDelete(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		customerID := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "John Doe", "11144477735", "CPF", "11987654321", "john@example.com")
		vehicleID := testsupport.ThereIsAVehicle(t, setup.Container.DB, uuid.Nil, customerID, "ABC1234", "Toyota", "Corolla", 2020)

		resp := deleteVehicle(t, setup, vehicleID)
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		resp2 := getVehicle(t, setup, vehicleID)
		require.Equal(t, fiber.StatusNotFound, resp2.StatusCode)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := deleteVehicle(t, setup, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestVehicleGetByID(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		customerID := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "John Doe", "11144477735", "CPF", "11987654321", "john@example.com")
		vehicleID := testsupport.ThereIsAVehicle(t, setup.Container.DB, uuid.Nil, customerID, "ABC1234", "Toyota", "Corolla", 2020)

		resp := getVehicle(t, setup, vehicleID)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result vehicleJSON
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, vehicleID, result.ID)
		assert.Equal(t, customerID, result.CustomerID)
		assert.Equal(t, "ABC1234", result.Plate)
		assert.Equal(t, "Toyota", result.Brand)
		assert.Equal(t, "Corolla", result.Model)
		assert.Equal(t, 2020, result.Year)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := getVehicle(t, setup, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestVehicleList(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		customerID := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "John Doe", "11144477735", "CPF", "11987654321", "john@example.com")
		testsupport.ThereIsAVehicle(t, setup.Container.DB, uuid.Nil, customerID, "ABC1234", "Toyota", "Corolla", 2020)
		testsupport.ThereIsAVehicle(t, setup.Container.DB, uuid.Nil, customerID, "XYZ5678", "Honda", "Civic", 2019)

		resp := getVehicles(t, setup)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result struct {
			Data []vehicleJSON `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Len(t, result.Data, 5)
	})
}

func TestVehicleListByCustomer(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		customerID1 := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "John Doe", "11144477735", "CPF", "11987654321", "john@example.com")
		customerID2 := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "Jane Doe", "52998224725", "CPF", "11987654322", "jane@example.com")

		testsupport.ThereIsAVehicle(t, setup.Container.DB, uuid.Nil, customerID1, "ABC1234", "Toyota", "Corolla", 2020)
		testsupport.ThereIsAVehicle(t, setup.Container.DB, uuid.Nil, customerID1, "XYZ5678", "Honda", "Civic", 2019)

		testsupport.ThereIsAVehicle(t, setup.Container.DB, uuid.Nil, customerID2, "DEF9012", "Ford", "Focus", 2021)

		resp := getVehiclesByCustomer(t, setup, customerID1)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result struct {
			Data []vehicleJSON `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Len(t, result.Data, 2)
		for _, vehicle := range result.Data {
			assert.Equal(t, customerID1, vehicle.CustomerID)
		}
	})
}

// -------- Helper Functions --------

type createBody struct {
	CustomerID uuid.UUID `json:"customer_id"`
	Plate      string    `json:"plate"`
	Brand      string    `json:"brand"`
	Model      string    `json:"model"`
	Year       int       `json:"year"`
}

type updateBody struct {
	Plate *string `json:"plate"`
	Brand *string `json:"brand"`
	Model *string `json:"model"`
	Year  *int    `json:"year"`
}

type vehicleJSON struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Plate      string    `json:"plate"`
	Brand      string    `json:"brand"`
	Model      string    `json:"model"`
	Year       int       `json:"year"`
}

func postCreateVehicle(t *testing.T, setup *testsupport.SetupConfig, body createBody) *http.Response {
	t.Helper()
	return doJSON(t, setup, "POST", "/admin/vehicles", body)
}

func patchUpdateVehicle(t *testing.T, setup *testsupport.SetupConfig, vehicleID uuid.UUID, body updateBody) *http.Response {
	t.Helper()
	id := uuidAsString(t, vehicleID)
	return doJSON(t, setup, "PATCH", "/admin/vehicles/"+id, body)
}

func deleteVehicle(t *testing.T, setup *testsupport.SetupConfig, vehicleID uuid.UUID) *http.Response {
	t.Helper()
	id := uuidAsString(t, vehicleID)
	return doJSON(t, setup, "DELETE", "/admin/vehicles/"+id, nil)
}

func getVehicle(t *testing.T, setup *testsupport.SetupConfig, vehicleID uuid.UUID) *http.Response {
	t.Helper()
	id := uuidAsString(t, vehicleID)
	return doJSON(t, setup, "GET", "/admin/vehicles/"+id, nil)
}

func getVehicles(t *testing.T, setup *testsupport.SetupConfig) *http.Response {
	t.Helper()
	return doJSON(t, setup, "GET", "/admin/vehicles", nil)
}

func getVehiclesByCustomer(t *testing.T, setup *testsupport.SetupConfig, customerID uuid.UUID) *http.Response {
	t.Helper()
	id := uuidAsString(t, customerID)
	return doJSON(t, setup, "GET", "/admin/customers/"+id+"/vehicles", nil)
}

func doJSON(t *testing.T, setup *testsupport.SetupConfig, method, path string, payload any) *http.Response {
	t.Helper()

	var bodyReader *bytes.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		require.NoError(t, err, "marshal payload")
		bodyReader = bytes.NewReader(b)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	testauth.AddAuthHeader(req, setup.AuthToken)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := setup.Container.FiberApp.Test(req, -1)
	require.NoError(t, err, "fiber app.Test")

	t.Cleanup(func() { _ = resp.Body.Close() })
	return resp
}

func uuidAsString(t *testing.T, id uuid.UUID) string {
	t.Helper()
	if id == uuid.Nil {
		return "invalid"
	}
	return id.String()
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
