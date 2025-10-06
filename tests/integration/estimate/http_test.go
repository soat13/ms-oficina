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
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/estimate"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/repairorder"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	repairorderShared "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
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

func ensureSetup(t *testing.T) *testsupport.SetupConfig {
	t.Helper()

	return testsupport.SetupHTTP(t, func(app *fiber.App, container *bootstrap.Container) {
		estimate.SetupDefault(container)
		repairorder.SetupDefault(container)
	})
}

func TestApprove(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := ThereIsARepairOrderInDiagnostics(t, setup.Container)
		require.Equal(t, fiber.StatusCreated, postCreateEstimate(t, repairOrderID, setup.Container).StatusCode)

		estimateID := getEstimateIDByRepairOrder(t, setup.Container, repairOrderID)

		resp := postApproveEstimate(t, setup.Container, estimateID)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
		expectEstimateStatus(t, setup.Container, repairOrderID, domain.StatusAwaitingStock)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusAwaitingApproval)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		resp := postApproveEstimateRaw(t, setup.Container, "invalid")
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postApproveEstimate(t, setup.Container, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Twice Should Fail", func(t *testing.T) {
		roID := ThereIsARepairOrderInDiagnostics(t, setup.Container)
		require.Equal(t, fiber.StatusCreated, postCreateEstimate(t, roID, setup.Container).StatusCode)
		estimateID := getEstimateIDByRepairOrder(t, setup.Container, roID)

		require.Equal(t, fiber.StatusOK, postApproveEstimate(t, setup.Container, estimateID).StatusCode)
		resp := postApproveEstimate(t, setup.Container, estimateID)

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestCreateFromRepairOrder(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := ThereIsARepairOrderInDiagnostics(t, setup.Container)

		resp := postCreateEstimate(t, repairOrderID, setup.Container)

		require.Equal(t, fiber.StatusCreated, resp.StatusCode)
		expectEstimateItemCount(t, setup.Container, repairOrderID, 2)
		expectEstimateStatus(t, setup.Container, repairOrderID, domain.StatusAwaitingApproval)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusAwaitingApproval)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		resp := postCreateEstimate(t, uuid.Nil, setup.Container)

		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		expectEstimateItemCount(t, setup.Container, uuid.Nil, 0)
	})

	t.Run("Invalid Status", func(t *testing.T) {
		repairOrderID := ThereIsARepairOrderReceived(t, setup.Container)

		resp := postCreateEstimate(t, repairOrderID, setup.Container)

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
		expectEstimateItemCount(t, setup.Container, repairOrderID, 0)
	})
}

func postCreateEstimate(t *testing.T, repairOrderID uuid.UUID, container *bootstrap.Container) *http.Response {
	t.Helper()

	payload := buildEstimatePayload(1, 1)

	var id string
	if repairOrderID == uuid.Nil {
		id = "invalid"
	} else {
		id = repairOrderID.String()
	}

	return DoJSON(t, container.FiberApp, "POST", "/admin/repair-orders/"+id+"/estimate", payload)
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

func ThereIsARepairOrderInDiagnostics(t *testing.T, container *bootstrap.Container) uuid.UUID {
	t.Helper()
	return testsupport.ThereIsARepairOrderInDiagnostics(t, container.DB)
}

func ThereIsARepairOrderReceived(t *testing.T, container *bootstrap.Container) uuid.UUID {
	t.Helper()
	return testsupport.ThereIsAReceivedRepairOrder(t, container.DB)
}

func expectEstimateItemCount(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID, expected int) {
	t.Helper()
	ctx := context.Background()

	var count int
	require.NoError(t,
		container.DB.NewRaw(`
			SELECT COUNT(*)
			FROM estimate_items ei
			JOIN estimates e ON e.id = ei.estimate_id
			WHERE e.repair_id = ?
		`, repairOrderID).Scan(ctx, &count),
	)
	assert.Equal(t, expected, count)
}

func expectEstimateStatus(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID, expected domain.Status) {
	t.Helper()
	ctx := context.Background()

	var status string
	require.NoError(t,
		container.DB.NewRaw(`
			SELECT e.status
			FROM estimates e
			JOIN repair_orders ro ON ro.id = e.repair_id
			WHERE ro.id = ?
		`, repairOrderID).Scan(ctx, &status),
	)

	assert.Equal(t, string(expected), status)
}

func expectRepairOrderStatus(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID, expected repairorderShared.Status) {
	t.Helper()
	ctx := context.Background()

	var status string
	require.NoError(t,
		container.DB.NewRaw(`
			SELECT status
			FROM repair_orders
			WHERE id = ?
		`, repairOrderID).Scan(ctx, &status),
	)

	assert.Equal(t, string(expected), status)
}

func postApproveEstimate(t *testing.T, container *bootstrap.Container, estimateID uuid.UUID) *http.Response {
	t.Helper()
	return DoJSON(t, container.FiberApp, "POST", "/admin/estimates/"+estimateID.String()+"/approve", nil)
}

func postApproveEstimateRaw(t *testing.T, container *bootstrap.Container, id string) *http.Response {
	t.Helper()
	return DoJSON(t, container.FiberApp, "POST", "/admin/estimates/"+id+"/approve", nil)
}

func getEstimateIDByRepairOrder(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	var estimateID uuid.UUID
	require.NoError(t,
		container.DB.NewRaw(`
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
