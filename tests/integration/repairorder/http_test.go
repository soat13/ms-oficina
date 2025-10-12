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
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/repairorder"
	repairorderShared "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/tests/testsupport"
	testauth "github.com/soat13/fase-1-oficina/tests/testsupport/auth"
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
	})
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

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
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusFinished)
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
		repairOrderID := testsupport.ThereIsAFinishedRepairOrder(t, setup.Container.DB)

		resp := postReleaseVehicle(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusReleased)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postReleaseVehicle(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
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
