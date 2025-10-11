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

func expectRepairOrderExists(t *testing.T, container *bootstrap.Container, customerID, vehicleID uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	var count int
	require.NoError(t,
		container.DB.
			NewRaw(`SELECT COUNT(1) FROM repair_orders WHERE customer_id = ? AND vehicle_id = ?`, customerID, vehicleID).
			Scan(ctx, &count),
	)

	require.Equal(t, 1, count, "repair order should be created for the given customer and vehicle")
}
