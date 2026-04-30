package repairorder

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/soat13/ms-oficina/internal/bootstrap"
	"github.com/soat13/ms-oficina/internal/bootstrap/estimate"
	"github.com/soat13/ms-oficina/internal/bootstrap/repairorder"
	estimatDomain "github.com/soat13/ms-oficina/internal/estimate/domain"
	repairorderShared "github.com/soat13/ms-oficina/internal/shared/repairorder"
	"github.com/soat13/ms-oficina/tests/testsupport"
	testauth "github.com/soat13/ms-oficina/tests/testsupport/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Setup
// -----------------------------------------------------------------------------

func ensureSetup(t *testing.T) *testsupport.SetupConfig {
	t.Helper()

	return testsupport.SetupHTTP(t, func(app *fiber.App, container *bootstrap.Container) {
		repairorder.SetupDefault(container)
		estimate.SetupDefault(container)
	})
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func TestRepairOrderCancel(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsAnApprovedRepairOrder(t, setup.Container.DB)

		resp := postCancel(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusCanceled)
	})

	t.Run("Idempotent - cancel twice returns 204", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsAnApprovedRepairOrder(t, setup.Container.DB)

		resp1 := postCancel(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		require.Equal(t, fiber.StatusNoContent, resp1.StatusCode)

		resp2 := postCancel(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		require.Equal(t, fiber.StatusNoContent, resp2.StatusCode)

		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusCanceled)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postCancel(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Cannot cancel when released", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusReleased)

		resp := postCancel(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		require.Equal(t, fiber.StatusConflict, resp.StatusCode)

		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusReleased)
	})

	t.Run("Canceling RO also cancels its Estimate", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsAnApprovedRepairOrder(t, setup.Container.DB)
		testsupport.ThereIsAnEstimateForRepairOrder(t, setup.Container, repairOrderID)

		resp := postCancel(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusCanceled)
		expectEstimateStatusByRO(t, setup.Container, repairOrderID, estimatDomain.StatusCanceled)
	})
}

func TestRepairOrderCreate(t *testing.T) {
	setup := ensureSetup(t)

	customerID := testsupport.ThereIsACustomerWithDocument(t, setup.Container.DB, uuid.Nil, "João Silva", "11144477735", "CPF", "11987654321", "joao@example.com")
	vehicleID := testsupport.ThereIsAVehicle(t, setup.Container.DB, uuid.Nil, customerID, "ABB1234", "Toyota", "Corolla", 2020)

	t.Run("Success", func(t *testing.T) {
		resp := postCreateRepairOrder(t, setup.Container.FiberApp, setup.AuthToken, customerID, vehicleID)

		require.Equal(t, fiber.StatusCreated, resp.StatusCode)
		expectRepairOrderExists(t, setup.Container, customerID, vehicleID)
	})

	t.Run("VehicleNotFound", func(t *testing.T) {
		unknownVehicleID := uuid.New()

		resp := postCreateRepairOrder(t, setup.Container.FiberApp, setup.AuthToken, customerID, unknownVehicleID)

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("CustomerNotFound", func(t *testing.T) {
		unknownCustomerID := uuid.New()

		resp := postCreateRepairOrder(t, setup.Container.FiberApp, setup.AuthToken, unknownCustomerID, vehicleID)

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRepairOrderStartDiagnostics(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsAReceivedRepairOrder(t, setup.Container.DB)

		resp := postStartDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusInDiagnostics)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postStartDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/admin/repair-orders/invalid-uuid/start-diagnostics", nil)
		testauth.AddAuthHeader(req, setup.AuthToken)
		resp, err := setup.Container.FiberApp.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestRepairOrderFinishDiagnostics(t *testing.T) {
	setup := ensureSetup(t)

	productID1 := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Oil Filter", 5000, 10)
	productID2 := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Brake Pad", 15000, 5)
	serviceID1 := testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Engine Oil Change", 12000)
	serviceID2 := testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Brake Check", 8000)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		products := []map[string]interface{}{
			{"id": productID1.String(), "quantity": 2},
			{"id": productID2.String(), "quantity": 1},
		}
		services := []map[string]interface{}{
			{"id": serviceID1.String(), "quantity": 1},
			{"id": serviceID2.String(), "quantity": 1},
		}

		resp := postFinishDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, products, services)

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusAwaitingApproval)
	})

	t.Run("Not Found", func(t *testing.T) {
		products := []map[string]interface{}{
			{"id": productID1.String(), "quantity": 1},
		}
		services := []map[string]interface{}{
			{"id": serviceID1.String(), "quantity": 1},
		}

		resp := postFinishDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, uuid.New(), products, services)
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Empty Products", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		products := []map[string]interface{}{}
		services := []map[string]interface{}{
			{"id": serviceID1.String(), "quantity": 1},
		}

		resp := postFinishDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, products, services)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Empty Services", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		products := []map[string]interface{}{
			{"id": productID1.String(), "quantity": 1},
		}
		services := []map[string]interface{}{}

		resp := postFinishDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, products, services)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Invalid Quantity - Zero", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		products := []map[string]interface{}{
			{"id": productID1.String(), "quantity": 0},
		}
		services := []map[string]interface{}{
			{"id": serviceID1.String(), "quantity": 1},
		}

		resp := postFinishDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, products, services)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Invalid Quantity - Negative", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		products := []map[string]interface{}{
			{"id": productID1.String(), "quantity": -1},
		}
		services := []map[string]interface{}{
			{"id": serviceID1.String(), "quantity": 1},
		}

		resp := postFinishDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, products, services)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Product Not Found", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		products := []map[string]interface{}{
			{"id": uuid.New().String(), "quantity": 1},
		}
		services := []map[string]interface{}{
			{"id": serviceID1.String(), "quantity": 1},
		}

		resp := postFinishDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, products, services)
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Service Not Found", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		products := []map[string]interface{}{
			{"id": productID1.String(), "quantity": 1},
		}
		services := []map[string]interface{}{
			{"id": uuid.New().String(), "quantity": 1},
		}

		resp := postFinishDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, products, services)
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Insufficient Stock", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		products := []map[string]interface{}{
			{"id": productID1.String(), "quantity": 100},
		}
		services := []map[string]interface{}{
			{"id": serviceID1.String(), "quantity": 1},
		}

		resp := postFinishDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, products, services)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Invalid UUID in Product", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		products := []map[string]interface{}{
			{"id": "invalid-uuid", "quantity": 1},
		}
		services := []map[string]interface{}{
			{"id": serviceID1.String(), "quantity": 1},
		}

		resp := postFinishDiagnostics(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, products, services)
		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		req := httptest.NewRequest("POST", "/admin/repair-orders/"+repairOrderID.String()+"/finish-diagnostics", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		testauth.AddAuthHeader(req, setup.AuthToken)

		resp, err := setup.Container.FiberApp.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestRepairOrderStartExecution(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsAnApprovedRepairOrder(t, setup.Container.DB)

		resp := postStartExecution(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusInExecution)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postStartExecution(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestRepairOrderFinishExecution(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInExecution(t, setup.Container.DB)

		resp := postFinishExecution(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusPaymentRequested)
		expectRepairOrderHasExecutionTime(t, setup.Container, repairOrderID)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postFinishExecution(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestRepairOrderReleaseVehicle(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsAPaymentSucceededRepairOrder(t, setup.Container.DB)

		resp := postReleaseVehicle(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusReleased)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postReleaseVehicle(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestRepairOrderGet(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsAReceivedRepairOrder(t, setup.Container.DB)

		resp := getRepairOrder(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, repairOrderID.String(), result["id"])
		require.Equal(t, "received", result["status"])
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := getRepairOrder(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/repair-orders/invalid-uuid", nil)
		testauth.AddAuthHeader(req, setup.AuthToken)
		resp, err := setup.Container.FiberApp.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestRepairOrderList(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success - Returns List", func(t *testing.T) {
		cleanAllRepairOrders(t, setup.Container)

		testsupport.ThereIsAnApprovedRepairOrder(t, setup.Container.DB)
		testsupport.ThereIsARepairOrderInExecution(t, setup.Container.DB)
		testsupport.ThereIsAnAwaitingApproveRepairOrder(t, setup.Container.DB)
		testsupport.ThereIsAFinishedRepairOrder(t, setup.Container.DB)
		testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)
		testsupport.ThereIsAReleasedRepairOrder(t, setup.Container.DB)
		testsupport.ThereIsAReceivedRepairOrder(t, setup.Container.DB)
		testsupport.ThereIsACanceledRepairOrder(t, setup.Container.DB)
		testsupport.ThereIsARepairOrderInDiagnosticsFinished(t, setup.Container.DB)

		resp := listRepairOrders(t, setup.Container.FiberApp, setup.AuthToken, "", "")
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.Contains(t, result, "data")

		data, ok := result["data"].([]interface{})
		require.True(t, ok)
		assert.Len(t, data, 6)

		getStatusFn := func(i int) string {
			item, _ := data[i].(map[string]interface{})
			status, _ := item["status"].(string)

			return status
		}

		assert.Equal(t, string(repairorderShared.StatusInExecution), getStatusFn(0))
		assert.Equal(t, string(repairorderShared.StatusApproved), getStatusFn(1))
		assert.Equal(t, string(repairorderShared.StatusAwaitingApproval), getStatusFn(2))
		assert.Equal(t, string(repairorderShared.StatusDiagnosticsFinished), getStatusFn(3))
		assert.Equal(t, string(repairorderShared.StatusInDiagnostics), getStatusFn(4))
		assert.Equal(t, string(repairorderShared.StatusReceived), getStatusFn(5))
	})

	t.Run("Success - With Pagination", func(t *testing.T) {
		resp := listRepairOrders(t, setup.Container.FiberApp, setup.AuthToken, "2", "0")

		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.Contains(t, result, "data")

		data, ok := result["data"].([]interface{})
		require.True(t, ok)
		require.LessOrEqual(t, len(data), 2)
	})

	t.Run("Success - Empty List", func(t *testing.T) {
		cleanAllRepairOrders(t, setup.Container)

		resp := listRepairOrders(t, setup.Container.FiberApp, setup.AuthToken, "", "")

		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.Contains(t, result, "data")

		data, ok := result["data"].([]interface{})
		require.True(t, ok)
		require.Equal(t, 0, len(data))
	})
}

func TestRepairOrderGetAverageExecutionTime(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success - With Finished Orders", func(t *testing.T) {
		ro1 := testsupport.ThereIsARepairOrderInExecution(t, setup.Container.DB)
		ro2 := testsupport.ThereIsARepairOrderInExecution(t, setup.Container.DB)
		ro3 := testsupport.ThereIsARepairOrderInExecution(t, setup.Container.DB)

		setExecutionTimeAndFinish(t, setup.Container, ro1, 60)
		setExecutionTimeAndFinish(t, setup.Container, ro2, 120)
		setExecutionTimeAndFinish(t, setup.Container, ro3, 180)

		resp := getAverageExecutionTime(t, setup.Container.FiberApp, setup.AuthToken)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result struct {
			AverageMinutes *float64 `json:"AverageMinutes"`
		}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		require.NotNil(t, result.AverageMinutes)
		assert.InDelta(t, 120.0, *result.AverageMinutes, 0.1)
	})

	t.Run("Success - No Finished Orders", func(t *testing.T) {
		cleanFinishedRepairOrders(t, setup.Container)

		resp := getAverageExecutionTime(t, setup.Container.FiberApp, setup.AuthToken)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result struct {
			AverageMinutes *float64 `json:"AverageMinutes"`
		}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		require.NotNil(t, result.AverageMinutes)
		assert.Equal(t, 0.0, *result.AverageMinutes)
	})

	t.Run("Success - Only Finished Orders Count", func(t *testing.T) {
		cleanFinishedRepairOrders(t, setup.Container)

		roInExecution := testsupport.ThereIsARepairOrderInExecution(t, setup.Container.DB)
		roApproved := testsupport.ThereIsAnApprovedRepairOrder(t, setup.Container.DB)
		roFinished1 := testsupport.ThereIsARepairOrderInExecution(t, setup.Container.DB)
		roFinished2 := testsupport.ThereIsARepairOrderInExecution(t, setup.Container.DB)

		setExecutionTime(t, setup.Container, roInExecution, 100)
		setExecutionTime(t, setup.Container, roApproved, 200)
		setExecutionTimeAndFinish(t, setup.Container, roFinished1, 60)
		setExecutionTimeAndFinish(t, setup.Container, roFinished2, 90)

		resp := getAverageExecutionTime(t, setup.Container.FiberApp, setup.AuthToken)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result struct {
			AverageMinutes *float64 `json:"AverageMinutes"`
		}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		require.NotNil(t, result.AverageMinutes)
		assert.InDelta(t, 75.0, *result.AverageMinutes, 0.1)
	})
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func expectRepairOrderStatus(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID, expected repairorderShared.Status) {
	t.Helper()
	ctx := context.Background()

	var status string
	require.NoError(t,
		container.DB.NewRaw(`SELECT status FROM repair_orders WHERE id = ?`, repairOrderID).Scan(ctx, &status),
	)

	assert.Equal(t, string(expected), status)
}

func expectRepairOrderHasExecutionTime(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	var executionTimeMinutes *int64
	require.NoError(t,
		container.DB.NewRaw(`SELECT execution_time_minutes FROM repair_orders WHERE id = ?`, repairOrderID).Scan(ctx, &executionTimeMinutes),
	)

	require.NotNil(t, executionTimeMinutes)
}

func postStartExecution(t *testing.T, fiberApp *fiber.App, token string, repairOrderID uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("POST", "/admin/repair-orders/"+repairOrderID.String()+"/start-execution", nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func postFinishExecution(t *testing.T, fiberApp *fiber.App, token string, repairOrderID uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("POST", "/admin/repair-orders/"+repairOrderID.String()+"/finish-execution", nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func postReleaseVehicle(t *testing.T, fiberApp *fiber.App, token string, repairOrderID uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("POST", "/admin/repair-orders/"+repairOrderID.String()+"/release-vehicle", nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func postCreateRepairOrder(t *testing.T, app *fiber.App, token string, customerID, vehicleID uuid.UUID) *http.Response {
	t.Helper()

	body := map[string]any{
		"customer_id": customerID,
		"vehicle_id":  vehicleID,
	}
	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/admin/repair-orders", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	testauth.AddAuthHeader(req, token)

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getAverageExecutionTime(t *testing.T, fiberApp *fiber.App, token string) *http.Response {
	t.Helper()
	req := httptest.NewRequest("GET", "/admin/repair-orders/average-execution-time", nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func expectRepairOrderExists(t *testing.T, container *bootstrap.Container, customerID, vehicleID uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	var count int
	require.NoError(t,
		container.DB.
			NewRaw(`SELECT COUNT(1) FROM repair_orders WHERE customer_id = ? AND vehicle_id = ?`, customerID, vehicleID).
			Scan(ctx, &count),
	)

	require.Equal(t, 1, count)
}

func setExecutionTime(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID, minutes int64) {
	t.Helper()
	ctx := context.Background()

	_, err := container.DB.NewRaw(`
		UPDATE repair_orders 
		SET execution_time_minutes = ?
		WHERE id = ?
	`, minutes, repairOrderID).Exec(ctx)
	require.NoError(t, err)
}

func setExecutionTimeAndFinish(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID, minutes int64) {
	t.Helper()
	ctx := context.Background()

	_, err := container.DB.NewRaw(`
		UPDATE repair_orders 
		SET execution_time_minutes = ?, status = ?
		WHERE id = ?
	`, minutes, repairorderShared.StatusFinished, repairOrderID).Exec(ctx)
	require.NoError(t, err)
}

func cleanFinishedRepairOrders(t *testing.T, container *bootstrap.Container) {
	t.Helper()
	ctx := context.Background()

	_, err := container.DB.NewRaw(`
		DELETE FROM repair_orders 
		WHERE status = ?
	`, repairorderShared.StatusFinished).Exec(ctx)
	require.NoError(t, err)
}

func postCancel(t *testing.T, fiberApp *fiber.App, token string, repairOrderID uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("POST", "/admin/repair-orders/"+repairOrderID.String()+"/cancel", nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func expectEstimateStatusByRO(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID, expected estimatDomain.Status) {
	t.Helper()
	ctx := context.Background()

	var status string
	require.NoError(t,
		container.DB.NewRaw(`
			SELECT status
			FROM estimates
			WHERE repair_order_id = ?
			ORDER BY created_at DESC
			LIMIT 1
		`, repairOrderID).Scan(ctx, &status),
	)

	assert.Equal(t, string(expected), status)
}

func postStartDiagnostics(t *testing.T, fiberApp *fiber.App, token string, repairOrderID uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("POST", "/admin/repair-orders/"+repairOrderID.String()+"/start-diagnostics", nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func postFinishDiagnostics(t *testing.T, fiberApp *fiber.App, token string, repairOrderID uuid.UUID, products, services []map[string]interface{}) *http.Response {
	t.Helper()

	body := map[string]any{
		"products": products,
		"services": services,
	}
	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/admin/repair-orders/"+repairOrderID.String()+"/finish-diagnostics", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	testauth.AddAuthHeader(req, token)

	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getRepairOrder(t *testing.T, fiberApp *fiber.App, token string, repairOrderID uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("GET", "/admin/repair-orders/"+repairOrderID.String(), nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func listRepairOrders(t *testing.T, fiberApp *fiber.App, token string, limit, offset string) *http.Response {
	t.Helper()
	url := "/admin/repair-orders"
	if limit != "" || offset != "" {
		url += "?"
		if limit != "" {
			url += "limit=" + limit
		}
		if offset != "" {
			if limit != "" {
				url += "&"
			}
			url += "offset=" + offset
		}
	}
	req := httptest.NewRequest("GET", url, nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func cleanAllRepairOrders(t *testing.T, container *bootstrap.Container) {
	t.Helper()
	ctx := context.Background()

	_, err := container.DB.NewRaw(`DELETE FROM repair_orders`).Exec(ctx)
	require.NoError(t, err)
}
