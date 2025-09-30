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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/tests/testsupport"
)

func TestRepairOrderStartExecution(t *testing.T) {
	ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsAnApprovedRepairOrder(t, env.db)

		resp := postStartExecution(t, repairOrderID)

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		expectRepairOrderStatus(t, repairOrderID, repairorder.StatusInExecution)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postStartExecution(t, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func expectRepairOrderStatus(t *testing.T, repairOrderID uuid.UUID, expected repairorder.Status) {
	t.Helper()
	ctx := context.Background()

	var status string
	require.NoError(t,
		env.db.NewRaw(`SELECT status FROM repair_orders WHERE id = ?`, repairOrderID).Scan(ctx, &status),
	)

	assert.Equal(t, string(expected), status)
}

func postStartExecution(t *testing.T, repairOrderID uuid.UUID) *http.Response {
	t.Helper()

	var id string
	if repairOrderID == uuid.Nil {
		id = "invalid"
	} else {
		id = repairOrderID.String()
	}

	return DoJSON(t, env.app, "POST", "/admin/repair-orders/"+id+"/start-execution", nil)
}

func postFinishExecution(t *testing.T, repairOrderID uuid.UUID) *http.Response {
	t.Helper()

	var id string
	if repairOrderID == uuid.Nil {
		id = "invalid"
	} else {
		id = repairOrderID.String()
	}

	return DoJSON(t, env.app, "POST", "/admin/repair-orders/"+id+"/finish-execution", nil)
}

func DoJSON(t *testing.T, app *fiber.App, method, path string, payload any) *http.Response {
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

	resp, err := app.Test(req, -1)
	require.NoError(t, err, "fiber app.Test")

	t.Cleanup(func() { _ = resp.Body.Close() })
	return resp
}
