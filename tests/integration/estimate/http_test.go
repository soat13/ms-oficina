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
		estimate.SetupDefault(container)
		repairorder.SetupDefault(container)
	})
}

// -----------------------------------------------------------------------------
// DTOs
// -----------------------------------------------------------------------------

type estimateLine struct {
	ID       uuid.UUID `json:"id"`
	Quantity int       `json:"quantity"`
}

type estimateBody struct {
	Products []estimateLine `json:"products"`
	Services []estimateLine `json:"services"`
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------
func TestReject(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)
		require.Equal(t, fiber.StatusCreated, postCreateEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID).StatusCode)

		estimateID := getEstimateIDByRepairOrder(t, setup.Container, repairOrderID)

		resp := postRejectEstimate(t, setup.Container.FiberApp, setup.AuthToken, estimateID)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		expectEstimateStatus(t, setup.Container, repairOrderID, domain.StatusRejected)

		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusCanceled)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		resp := postRejectEstimateRaw(t, setup.Container.FiberApp, setup.AuthToken, "invalid")
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postRejectEstimate(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Twice Should Fail", func(t *testing.T) {
		roID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)
		require.Equal(t, fiber.StatusCreated, postCreateEstimate(t, setup.Container.FiberApp, setup.AuthToken, roID).StatusCode)

		estimateID := getEstimateIDByRepairOrder(t, setup.Container, roID)

		require.Equal(t, fiber.StatusOK, postRejectEstimate(t, setup.Container.FiberApp, setup.AuthToken, estimateID).StatusCode)
		resp := postRejectEstimate(t, setup.Container.FiberApp, setup.AuthToken, estimateID)

		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})
}

func TestApprove(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)
		require.Equal(t, fiber.StatusCreated, postCreateEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID).StatusCode)

		estimateID := getEstimateIDByRepairOrder(t, setup.Container, repairOrderID)

		resp := postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, estimateID)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
		expectEstimateStatus(t, setup.Container, repairOrderID, domain.StatusAwaitingStock)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusAwaitingApproval)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		resp := postApproveEstimateRaw(t, setup.Container.FiberApp, setup.AuthToken, "invalid")
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Not Found", func(t *testing.T) {
		resp := postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, uuid.New())
		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Twice Should Fail", func(t *testing.T) {
		roID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)
		require.Equal(t, fiber.StatusCreated, postCreateEstimate(t, setup.Container.FiberApp, setup.AuthToken, roID).StatusCode)
		estimateID := getEstimateIDByRepairOrder(t, setup.Container, roID)

		require.Equal(t, fiber.StatusOK, postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, estimateID).StatusCode)
		resp := postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, estimateID)

		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})
}

func TestCreateFromRepairOrder(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		resp := postCreateEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusCreated, resp.StatusCode)
		expectEstimateItemCount(t, setup.Container, repairOrderID, 2)
		expectEstimateStatus(t, setup.Container, repairOrderID, domain.StatusAwaitingApproval)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusAwaitingApproval)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		resp := postCreateEstimate(t, setup.Container.FiberApp, setup.AuthToken, uuid.Nil)

		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		expectEstimateItemCount(t, setup.Container, uuid.Nil, 0)
	})

	t.Run("Invalid Status", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsAReceivedRepairOrder(t, setup.Container.DB)

		resp := postCreateEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
		expectEstimateItemCount(t, setup.Container, repairOrderID, 0)
	})

	t.Run("No products and no services", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		body := estimateBody{
			Products: []estimateLine{},
			Services: []estimateLine{},
		}

		resp := postCreateEstimateWithBody(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, body)

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
		expectEstimateItemCount(t, setup.Container, repairOrderID, 0)
	})

	t.Run("Quantities equals to zero", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		body := estimateBody{
			Products: []estimateLine{
				{ID: testsupport.OilFilterID, Quantity: 0},
			},
			Services: nil,
		}

		resp := postCreateEstimateWithBody(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, body)

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
		expectEstimateItemCount(t, setup.Container, repairOrderID, 0)
	})

	t.Run("Insufficient stock", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderInDiagnostics(t, setup.Container.DB)

		body := estimateBody{
			Products: []estimateLine{
				{ID: testsupport.OilFilterID, Quantity: 300},
			},
			Services: nil,
		}

		resp := postCreateEstimateWithBody(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID, body)

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
		expectEstimateItemCount(t, setup.Container, repairOrderID, 0)
	})
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

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

func expectEstimateItemCount(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID, expected int) {
	t.Helper()
	ctx := context.Background()

	var count int
	require.NoError(t,
		container.DB.NewRaw(`
			SELECT COUNT(*)
			FROM estimate_items ei
			JOIN estimates e ON e.id = ei.estimate_id
			WHERE e.repair_order_id = ?
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
			JOIN repair_orders ro ON ro.id = e.repair_order_id
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
func postCreateEstimate(t *testing.T, fiberApp *fiber.App, token string, repairOrderID uuid.UUID) *http.Response {
	t.Helper()
	payload := buildEstimatePayload(1, 1)

	var id string
	if repairOrderID == uuid.Nil {
		id = "invalid"
	} else {
		id = repairOrderID.String()
	}

	bs, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/admin/repair-orders/"+id+"/estimate", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func postCreateEstimateWithBody(t *testing.T, fiberApp *fiber.App, token string, repairOrderID uuid.UUID, body estimateBody) *http.Response {
	t.Helper()

	var id string
	if repairOrderID == uuid.Nil {
		id = "invalid"
	} else {
		id = repairOrderID.String()
	}

	bs, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/admin/repair-orders/"+id+"/estimate", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func postApproveEstimate(t *testing.T, fiberApp *fiber.App, token string, estimateID uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("POST", "/admin/estimates/"+estimateID.String()+"/approve", nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func postApproveEstimateRaw(t *testing.T, fiberApp *fiber.App, token string, id string) *http.Response {
	t.Helper()
	req := httptest.NewRequest("POST", "/admin/estimates/"+id+"/approve", nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func postRejectEstimate(t *testing.T, fiberApp *fiber.App, token string, estimateID uuid.UUID) *http.Response {
	t.Helper()
	req := httptest.NewRequest("POST", "/admin/estimates/"+estimateID.String()+"/reject", nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func postRejectEstimateRaw(t *testing.T, fiberApp *fiber.App, token string, id string) *http.Response {
	t.Helper()
	req := httptest.NewRequest("POST", "/admin/estimates/"+id+"/reject", nil)
	testauth.AddAuthHeader(req, token)
	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func getEstimateIDByRepairOrder(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	var estimateID uuid.UUID
	require.NoError(t,
		container.DB.NewRaw(`
			SELECT e.id
			FROM estimates e
			WHERE e.repair_order_id = ?
			ORDER BY e.created_at DESC
			LIMIT 1
		`, repairOrderID).Scan(ctx, &estimateID),
	)
	require.NotEqual(t, uuid.Nil, estimateID)
	return estimateID
}
