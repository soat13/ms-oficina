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

func ensureSetup(t *testing.T) *testsupport.SetupConfig {
	t.Helper()

	return testsupport.SetupHTTP(t, func(app *fiber.App, container *bootstrap.Container) {
		repairorder.SetupDefault(container)
	})
}

func TestRepairOrderReleaseVehicle(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsAFinishedRepairOrder(t, setup.Container.DB)

		resp := postReleaseVehicle(t, setup, repairOrderID)

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusReleased)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postStartExecution(t, setup, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestRepairOrderFinishExecution(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInExecution(t, setup.Container.DB)

		resp := postFinishExecution(t, setup, repairOrderID)

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusFinished)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postStartExecution(t, setup, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestRepairOrderStartExecution(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsAnApprovedRepairOrder(t, setup.Container.DB)

		resp := postStartExecution(t, setup, repairOrderID)

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusInExecution)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postStartExecution(t, setup, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func expectRepairOrderStatus(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID, expected repairorderShared.Status) {
	t.Helper()
	ctx := context.Background()

	var status string
	require.NoError(t,
		container.DB.NewRaw(`SELECT status FROM repair_orders WHERE id = ?`, repairOrderID).Scan(ctx, &status),
	)

	assert.Equal(t, string(expected), status)
}

func postStartExecution(t *testing.T, setup *testsupport.SetupConfig, repairOrderID uuid.UUID) *http.Response {
	t.Helper()
	id := UuidAsString(t, repairOrderID)
	return DoJSON(t, setup.Container.FiberApp, setup.AuthToken, "POST", "/admin/repair-orders/"+id+"/start-execution", nil)
}

func postFinishExecution(t *testing.T, setup *testsupport.SetupConfig, repairOrderID uuid.UUID) *http.Response {
	t.Helper()
	id := UuidAsString(t, repairOrderID)
	return DoJSON(t, setup.Container.FiberApp, setup.AuthToken, "POST", "/admin/repair-orders/"+id+"/finish-execution", nil)
}

func postReleaseVehicle(t *testing.T, setup *testsupport.SetupConfig, repairOrderID uuid.UUID) *http.Response {
	t.Helper()
	id := UuidAsString(t, repairOrderID)
	return DoJSON(t, setup.Container.FiberApp, setup.AuthToken, "POST", "/admin/repair-orders/"+id+"/release-vehicle", nil)
}

func DoJSON(t *testing.T, app *fiber.App, token, method, path string, payload any) *http.Response {
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
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		testauth.AddAuthHeader(req, token)
	}

	resp, err := app.Test(req, -1)
	require.NoError(t, err, "fiber app.Test")

	t.Cleanup(func() { _ = resp.Body.Close() })
	return resp
}

func UuidAsString(t *testing.T, id uuid.UUID) string {
	t.Helper()
	if id == uuid.Nil {
		return "invalid"
	}

	return id.String()
}
