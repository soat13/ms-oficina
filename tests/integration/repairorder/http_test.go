package repairorder

import (
	"context"
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
