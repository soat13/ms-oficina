package estimate

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soat13/fase-1-oficina/tests/testsupport"
)

type estimateLine struct {
	ID       uuid.UUID `json:"id"`
	Quantity int       `json:"quantity"`
}

type estimateBody struct {
	Products []estimateLine `json:"products"`
	Services []estimateLine `json:"services"`
}

func Test_Estimate_Approve(t *testing.T) {
	ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		roID := ThereIsARepairOrderInDiagnostics(t)
		require.Equal(t, fiber.StatusCreated, postCreateEstimate(t, roID).StatusCode)

		estimateID := getEstimateIDByRepairOrder(t, roID)

		resp := postApproveEstimate(t, estimateID)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
		expectEstimateStatus(t, roID, domain.StatusAwaitingStock)
		expectRepairOrderStatus(t, roID, repairorder.StatusApproved)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		resp := postApproveEstimateRaw(t, "invalid")
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postApproveEstimate(t, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Twice Should Fail", func(t *testing.T) {
		roID := ThereIsARepairOrderInDiagnostics(t)
		require.Equal(t, fiber.StatusCreated, postCreateEstimate(t, roID).StatusCode)
		estimateID := getEstimateIDByRepairOrder(t, roID)

		require.Equal(t, fiber.StatusOK, postApproveEstimate(t, estimateID).StatusCode)
		resp := postApproveEstimate(t, estimateID)

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func Test_Estimate_Create_FromRepairOrder(t *testing.T) {
	ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := ThereIsARepairOrderInDiagnostics(t)

		resp := postCreateEstimate(t, repairOrderID)

		require.Equal(t, fiber.StatusCreated, resp.StatusCode)
		expectEstimateItemCount(t, repairOrderID, 2)
		expectEstimateStatus(t, repairOrderID, domain.StatusAwaitingApproval)
		expectRepairOrderStatus(t, repairOrderID, repairorder.StatusAwaitingApproval)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		resp := postCreateEstimate(t, uuid.Nil)

		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		expectEstimateItemCount(t, uuid.Nil, 0)
	})

	t.Run("Invalid Status", func(t *testing.T) {
		repairOrderID := ThereIsARepairOrderReceived(t)

		resp := postCreateEstimate(t, repairOrderID)

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
		expectEstimateItemCount(t, repairOrderID, 0)
	})
}

func postCreateEstimate(t *testing.T, repairOrderID uuid.UUID) *http.Response {
	t.Helper()

	payload := buildEstimatePayload(1, 1)

	var id string
	if repairOrderID == uuid.Nil {
		id = "invalid"
	} else {
		id = repairOrderID.String()
	}

	return DoJSON(t, env.app, "POST", "/admin/repair-orders/"+id+"/estimate", payload)
}

func buildEstimatePayload(productQty, serviceQty int) estimateBody {
	return estimateBody{
		Products: []estimateLine{
			{ID: testsupport.OilFilterID, Quantity: productQty},
		},
		Services: []estimateLine{
			{ID: testsupport.EngineOilChangeID, Quantity: serviceQty},
		},
	}
}

func ThereIsARepairOrderInDiagnostics(t *testing.T) uuid.UUID {
	t.Helper()
	return testsupport.ThereIsARepairOrderInDiagnostics(t, env.db)
}

func ThereIsARepairOrderReceived(t *testing.T) uuid.UUID {
	t.Helper()
	return testsupport.ThereIsAReceivedRepairOrder(t, env.db)
}

func expectEstimateItemCount(t *testing.T, repairOrderID uuid.UUID, expected int) {
	t.Helper()
	ctx := context.Background()

	var count int
	require.NoError(t,
		env.db.NewRaw(`
			SELECT COUNT(*)
			FROM estimate_items ei
			JOIN estimates e ON e.id = ei.estimate_id
			WHERE e.repair_id = ?
		`, repairOrderID).Scan(ctx, &count),
	)
	assert.Equal(t, expected, count)
}

func expectEstimateStatus(t *testing.T, repairOrderID uuid.UUID, expected domain.Status) {
	t.Helper()
	ctx := context.Background()

	var status string
	require.NoError(t,
		env.db.NewRaw(`
			SELECT e.status
			FROM estimates e
			JOIN repair_orders ro ON ro.id = e.repair_id
			WHERE ro.id = ?
		`, repairOrderID).Scan(ctx, &status),
	)

	assert.Equal(t, string(expected), status)
}

func expectRepairOrderStatus(t *testing.T, repairOrderID uuid.UUID, expected repairorder.Status) {
	t.Helper()
	ctx := context.Background()

	var status string
	require.NoError(t,
		env.db.NewRaw(`
			SELECT status
			FROM repair_orders
			WHERE id = ?
		`, repairOrderID).Scan(ctx, &status),
	)

	assert.Equal(t, string(expected), status)
}

func postApproveEstimate(t *testing.T, estimateID uuid.UUID) *http.Response {
	t.Helper()
	return DoJSON(t, env.app, "POST", "/admin/estimates/"+estimateID.String()+"/approve", nil)
}

func postApproveEstimateRaw(t *testing.T, id string) *http.Response {
	t.Helper()
	return DoJSON(t, env.app, "POST", "/admin/estimates/"+id+"/approve", nil)
}

func getEstimateIDByRepairOrder(t *testing.T, repairOrderID uuid.UUID) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	var estimateID uuid.UUID
	require.NoError(t,
		env.db.NewRaw(`
			SELECT e.id
			FROM estimates e
			WHERE e.repair_id = ?
			LIMIT 1
		`, repairOrderID).Scan(ctx, &estimateID),
	)
	require.NotEqual(t, uuid.Nil, estimateID)
	return estimateID
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
