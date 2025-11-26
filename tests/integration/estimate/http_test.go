package estimate

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/estimate"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/product"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/repairorder"
	estimateDomain "github.com/soat13/fase-1-oficina/internal/estimate/domain"
	repairorderShared "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/tests/testsupport"
	testauth "github.com/soat13/fase-1-oficina/tests/testsupport/auth"
)

// -----------------------------------------------------------------------------
// Setup
// -----------------------------------------------------------------------------

func ensureSetup(t *testing.T) *testsupport.SetupConfig {
	t.Helper()

	return testsupport.SetupHTTP(t, func(app *fiber.App, container *bootstrap.Container) {
		// Setup all modules to test complete event flows
		product.SetupDefault(container)     // Registers reduce_stock listener
		repairorder.SetupDefault(container) // Registers approval/cancel listeners
		estimate.SetupDefault(container)    // Registers estimate listeners
	})
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

// TestCreateEstimate - Create é executado APENAS via evento (OnDiagnosticsFinished listener)
// Testado em: tests/integration/repairorder/http_test.go -> TestFinishDiagnostics

func TestApproveEstimate(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := testsupport.ThereIsAnEstimateForRepairOrder(t, setup.Container, repairOrderID)

		resp := postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
		expectEstimateStatus(t, setup.Container, estimateID, estimateDomain.StatusApproved)
	})

	t.Run("Estimate Not Found", func(t *testing.T) {
		repairOrderID := uuid.New()

		resp := postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Cannot Approve Already Approved", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := testsupport.ThereIsAnEstimateForRepairOrder(t, setup.Container, repairOrderID)

		resp1 := postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		require.Equal(t, fiber.StatusOK, resp1.StatusCode)

		resp2 := postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusConflict, resp2.StatusCode)
		expectEstimateStatus(t, setup.Container, estimateID, estimateDomain.StatusApproved)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/admin/repair-orders/invalid-uuid/estimate/approve", nil)
		testauth.AddAuthHeader(req, setup.AuthToken)
		resp, err := setup.Container.FiberApp.Test(req, -1)

		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestRejectEstimate(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := testsupport.ThereIsAnEstimateForRepairOrder(t, setup.Container, repairOrderID)

		resp := postRejectEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
		expectEstimateStatus(t, setup.Container, estimateID, estimateDomain.StatusRejected)
	})

	t.Run("Estimate Not Found", func(t *testing.T) {
		repairOrderID := uuid.New()

		resp := postRejectEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Cannot Reject Already Rejected", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := testsupport.ThereIsAnEstimateForRepairOrder(t, setup.Container, repairOrderID)

		resp1 := postRejectEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		require.Equal(t, fiber.StatusOK, resp1.StatusCode)

		resp2 := postRejectEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)

		require.Equal(t, fiber.StatusConflict, resp2.StatusCode)
		expectEstimateStatus(t, setup.Container, estimateID, estimateDomain.StatusRejected)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/admin/repair-orders/invalid-uuid/estimate/reject", nil)
		testauth.AddAuthHeader(req, setup.AuthToken)
		resp, err := setup.Container.FiberApp.Test(req, -1)

		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestAddItemToEstimate(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success - Add Product Item", func(t *testing.T) {
		productID := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Oil Filter", 5000, 100)
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := testsupport.ThereIsAnEstimateForRepairOrder(t, setup.Container, repairOrderID)

		resp := postAddItem(t, setup.Container.FiberApp, setup.AuthToken, estimateID, productID, 2)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
		expectItemQuantity(t, setup.Container, estimateID, productID, 2)
	})

	t.Run("Success - Add Service Item", func(t *testing.T) {
		serviceID := testsupport.ThereIsAService(t, setup.Container.DB, uuid.Nil, "Alignment", 12000)
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := testsupport.ThereIsAnEstimateForRepairOrder(t, setup.Container, repairOrderID)

		resp := postAddItem(t, setup.Container.FiberApp, setup.AuthToken, estimateID, serviceID, 1)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
		expectItemQuantity(t, setup.Container, estimateID, serviceID, 1)
	})

	t.Run("Success - Add Same Item Twice Sums Quantity", func(t *testing.T) {
		productID := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Brake Pad", 15000, 50)
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := createEstimateWithProducts(t, setup.Container, repairOrderID, map[uuid.UUID]int{
			productID: 2,
		})

		resp := postAddItem(t, setup.Container.FiberApp, setup.AuthToken, estimateID, productID, 3)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		expectItemQuantity(t, setup.Container, estimateID, productID, 5) // 2 + 3
	})

	t.Run("Estimate Not Found", func(t *testing.T) {
		productID := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Oil Filter", 5000, 100)
		estimateID := uuid.New()

		resp := postAddItem(t, setup.Container.FiberApp, setup.AuthToken, estimateID, productID, 1)

		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Insufficient Stock", func(t *testing.T) {
		productID := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Low Stock Item", 5000, 5)
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := testsupport.ThereIsAnEstimateForRepairOrder(t, setup.Container, repairOrderID)

		resp := postAddItem(t, setup.Container.FiberApp, setup.AuthToken, estimateID, productID, 10)

		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})

	t.Run("Cannot Add After Approval", func(t *testing.T) {
		productID := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Oil Filter", 5000, 100)
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := testsupport.ThereIsAnEstimateForRepairOrder(t, setup.Container, repairOrderID)

		postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		time.Sleep(50 * time.Millisecond)

		resp := postAddItem(t, setup.Container.FiberApp, setup.AuthToken, estimateID, productID, 1)

		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/admin/estimates/invalid-uuid/items", nil)
		testauth.AddAuthHeader(req, setup.AuthToken)
		resp, err := setup.Container.FiberApp.Test(req, -1)

		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestRemoveItemFromEstimate(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Success", func(t *testing.T) {
		productID1 := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Oil Filter", 5000, 100)
		productID2 := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Brake Pad", 8000, 50)
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := createEstimateWithProducts(t, setup.Container, repairOrderID, map[uuid.UUID]int{
			productID1: 2,
			productID2: 1,
		})

		resp := deleteRemoveItem(t, setup.Container.FiberApp, setup.AuthToken, estimateID, productID1)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
		expectItemNotExists(t, setup.Container, estimateID, productID1)
	})

	t.Run("Estimate Not Found", func(t *testing.T) {
		estimateID := uuid.New()
		itemID := uuid.New()

		resp := deleteRemoveItem(t, setup.Container.FiberApp, setup.AuthToken, estimateID, itemID)

		require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Item Not Found", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := testsupport.ThereIsAnEstimateForRepairOrder(t, setup.Container, repairOrderID)
		itemID := uuid.New()

		resp := deleteRemoveItem(t, setup.Container.FiberApp, setup.AuthToken, estimateID, itemID)

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Cannot Remove Last Item", func(t *testing.T) {
		productID := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Oil Filter", 5000, 100)
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := createEstimateWithProducts(t, setup.Container, repairOrderID, map[uuid.UUID]int{
			productID: 1,
		})

		resp := deleteRemoveItem(t, setup.Container.FiberApp, setup.AuthToken, estimateID, productID)

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Cannot Remove After Approval", func(t *testing.T) {
		productID1 := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Oil Filter", 5000, 100)
		productID2 := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Brake Pad", 8000, 50)
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := createEstimateWithProducts(t, setup.Container, repairOrderID, map[uuid.UUID]int{
			productID1: 2,
			productID2: 1,
		})

		postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		time.Sleep(50 * time.Millisecond)

		resp := deleteRemoveItem(t, setup.Container.FiberApp, setup.AuthToken, estimateID, productID1)

		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})

	t.Run("Invalid Estimate UUID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/admin/estimates/invalid-uuid/items/"+uuid.New().String(), nil)
		testauth.AddAuthHeader(req, setup.AuthToken)
		resp, err := setup.Container.FiberApp.Test(req, -1)

		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Invalid Item UUID", func(t *testing.T) {
		estimateID := uuid.New()
		req := httptest.NewRequest("DELETE", "/admin/estimates/"+estimateID.String()+"/items/invalid-uuid", nil)
		testauth.AddAuthHeader(req, setup.AuthToken)
		resp, err := setup.Container.FiberApp.Test(req, -1)

		require.NoError(t, err)
		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

// -----------------------------------------------------------------------------
// Event Flow Tests - Testing complete event chains
// -----------------------------------------------------------------------------

func TestApproveEstimate_EventFlow(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Approve With Sufficient Stock - Triggers Stock Reduction and RO Approval", func(t *testing.T) {
		productID1 := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Oil Filter", 5000, 100)
		productID2 := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Brake Pad", 15000, 50)

		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)

		estimateID := createEstimateWithProducts(t, setup.Container, repairOrderID, map[uuid.UUID]int{
			productID1: 2,
			productID2: 1,
		})

		resp := postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		time.Sleep(100 * time.Millisecond)

		expectEstimateStatus(t, setup.Container, estimateID, estimateDomain.StatusApproved)

		expectProductStock(t, setup.Container, productID1, 98) // 100 - 2
		expectProductStock(t, setup.Container, productID2, 49) // 50 - 1

		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusApproved)
	})

	t.Run("Approve With Insufficient Stock - Triggers Stock Check and RO Cancellation", func(t *testing.T) {
		productID := testsupport.ThereIsAProduct(t, setup.Container.DB, uuid.Nil, "Low Stock Item", 5000, 5)

		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)

		estimateID := createEstimateWithProducts(t, setup.Container, repairOrderID, map[uuid.UUID]int{
			productID: 10, // Only 5 available
		})

		resp := postApproveEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		require.Contains(t, []int{fiber.StatusOK, fiber.StatusConflict}, resp.StatusCode, "Expected 200 or 409")

		time.Sleep(150 * time.Millisecond)

		var estimateStatus string
		ctx := context.Background()
		require.NoError(t, setup.Container.DB.NewRaw(`SELECT status FROM estimates WHERE id = ?`, estimateID).Scan(ctx, &estimateStatus))
		require.Contains(t, []string{string(estimateDomain.StatusApproved), string(estimateDomain.StatusCanceled)}, estimateStatus)

		expectProductStock(t, setup.Container, productID, 5) // Should remain unchanged

		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusCanceled)
	})
}

func TestRejectEstimate_EventFlow(t *testing.T) {
	setup := ensureSetup(t)

	t.Run("Reject Estimate - Triggers RepairOrder Cancellation", func(t *testing.T) {
		repairOrderID := testsupport.ThereIsARepairOrderWithStatus(t, setup.Container.DB, repairorderShared.StatusAwaitingApproval)
		estimateID := testsupport.ThereIsAnEstimateForRepairOrder(t, setup.Container, repairOrderID)

		resp := postRejectEstimate(t, setup.Container.FiberApp, setup.AuthToken, repairOrderID)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)

		time.Sleep(100 * time.Millisecond)

		expectEstimateStatus(t, setup.Container, estimateID, estimateDomain.StatusRejected)
		expectRepairOrderStatus(t, setup.Container, repairOrderID, repairorderShared.StatusCanceled)
	})
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func createEstimateWithProducts(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID, products map[uuid.UUID]int) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	estimateID := uuid.New()
	_, err := container.DB.NewRaw(`
		INSERT INTO estimates (id, repair_order_id, status)
		VALUES (?, ?, ?)
	`, estimateID, repairOrderID, "awaiting_approval").Exec(ctx)
	require.NoError(t, err)

	for productID, quantity := range products {
		_, err := container.DB.NewRaw(`
			INSERT INTO estimate_items (id, estimate_id, item_id, item_name, item_type, price, quantity)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, uuid.New(), estimateID, productID, "Product", "product", 1000, quantity).Exec(ctx)
		require.NoError(t, err)
	}

	return estimateID
}

func expectProductStock(t *testing.T, container *bootstrap.Container, productID uuid.UUID, expectedStock int) {
	t.Helper()
	ctx := context.Background()

	var actualStock int
	err := container.DB.NewRaw(`SELECT stock FROM products WHERE id = ?`, productID).Scan(ctx, &actualStock)
	require.NoError(t, err)
	assert.Equal(t, expectedStock, actualStock, "Product %s stock mismatch", productID)
}

func postApproveEstimate(t *testing.T, fiberApp *fiber.App, token string, repairOrderID uuid.UUID) *http.Response {
	t.Helper()

	req := httptest.NewRequest("POST", "/admin/repair-orders/"+repairOrderID.String()+"/estimate/approve", nil)
	testauth.AddAuthHeader(req, token)

	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func postRejectEstimate(t *testing.T, fiberApp *fiber.App, token string, repairOrderID uuid.UUID) *http.Response {
	t.Helper()

	req := httptest.NewRequest("POST", "/admin/repair-orders/"+repairOrderID.String()+"/estimate/reject", nil)
	testauth.AddAuthHeader(req, token)

	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func expectEstimateStatus(t *testing.T, container *bootstrap.Container, estimateID uuid.UUID, expectedStatus estimateDomain.Status) {
	t.Helper()
	ctx := context.Background()

	var status string
	require.NoError(t,
		container.DB.NewRaw(`SELECT status FROM estimates WHERE id = ?`, estimateID).
			Scan(ctx, &status),
	)
	assert.Equal(t, string(expectedStatus), status)
}

func expectRepairOrderStatus(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID, expectedStatus repairorderShared.Status) {
	t.Helper()
	ctx := context.Background()

	var status string
	require.NoError(t,
		container.DB.NewRaw(`SELECT status FROM repair_orders WHERE id = ?`, repairOrderID).
			Scan(ctx, &status),
	)
	assert.Equal(t, string(expectedStatus), status)
}

func postAddItem(t *testing.T, fiberApp *fiber.App, token string, estimateID uuid.UUID, itemID uuid.UUID, quantity int) *http.Response {
	t.Helper()

	body := fmt.Sprintf(`{"item_id":"%s","quantity":%d}`, itemID.String(), quantity)
	req := httptest.NewRequest("POST", "/admin/estimates/"+estimateID.String()+"/items", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	testauth.AddAuthHeader(req, token)

	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func deleteRemoveItem(t *testing.T, fiberApp *fiber.App, token string, estimateID uuid.UUID, itemID uuid.UUID) *http.Response {
	t.Helper()

	req := httptest.NewRequest("DELETE", "/admin/estimates/"+estimateID.String()+"/items/"+itemID.String(), nil)
	testauth.AddAuthHeader(req, token)

	resp, err := fiberApp.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func expectItemQuantity(t *testing.T, container *bootstrap.Container, estimateID uuid.UUID, itemID uuid.UUID, expectedQuantity int) {
	t.Helper()
	ctx := context.Background()

	var quantity int
	err := container.DB.NewRaw(`
		SELECT quantity FROM estimate_items 
		WHERE estimate_id = ? AND item_id = ?
	`, estimateID, itemID).Scan(ctx, &quantity)
	require.NoError(t, err)
	assert.Equal(t, expectedQuantity, quantity)
}

func expectItemNotExists(t *testing.T, container *bootstrap.Container, estimateID uuid.UUID, itemID uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	var count int
	err := container.DB.NewRaw(`
		SELECT COUNT(*) FROM estimate_items 
		WHERE estimate_id = ? AND item_id = ?
	`, estimateID, itemID).Scan(ctx, &count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "Item should not exist")
}

