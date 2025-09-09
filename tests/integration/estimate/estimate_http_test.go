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

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func Test_POST_CreateEstimateFromRepairOrder_OK(t *testing.T) {
	ensureSetup(t)

	// GIVEN
	repairOrderID := ThereIsARepairOrderInDiagnostics(t)

	// WHEN
	resp := postCreateEstimate(t, repairOrderID)

	// THEN
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
	estimateItemsCountForRepairOrder(t, repairOrderID, 2)
	// todo: improve this assert estimate and repair order status
}

func Test_POST_CreateEstimateFromRepairOrder_InvalidID(t *testing.T) {
	ensureSetup(t)

	// GIVEN
	repairOrderID := uuid.Nil

	// WHEN
	resp := postCreateEstimate(t, repairOrderID)

	// THEN
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	estimateItemsCountForRepairOrder(t, repairOrderID, 0)
}

func Test_POST_CreateEstimateFromRepairOrder_InvalidStatus(t *testing.T) {
	ensureSetup(t)

	// GIVEN
	repairOrderID := ThereIsARepairOrderReceived(t)

	// WHEN
	resp := postCreateEstimate(t, repairOrderID)

	// THEN
	require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	estimateItemsCountForRepairOrder(t, repairOrderID, 0)
}

// -----------------------------------------------------------------------------
// Helpers (sem precisar passar app/db em todo lugar)
// -----------------------------------------------------------------------------

func postCreateEstimate(t *testing.T, repairOrderID uuid.UUID) *http.Response {
	t.Helper()

	body := buildEstimateBody(t, 1, 1)

	var id string
	if repairOrderID == uuid.Nil {
		id = "invalid"
	} else {
		id = repairOrderID.String()
	}

	req := httptest.NewRequest(
		"POST",
		"/repair-orders/"+id+"/estimate",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := env.app.Test(req, -1)
	require.NoError(t, err)
	return resp
}

func buildEstimateBody(t *testing.T, productQty, serviceQty int) []byte {
	t.Helper()

	body := estimateBody{
		Products: []estimateLine{
			{ID: testsupport.OilFilterID, Quantity: productQty},
		},
		Services: []estimateLine{
			{ID: testsupport.EngineOilChangeID, Quantity: serviceQty},
		},
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)
	return b
}

// Givens wrappers usando env.db
func ThereIsARepairOrderInDiagnostics(t *testing.T) uuid.UUID {
	t.Helper()
	return testsupport.ThereIsARepairOrderInDiagnostics(t, env.db)
}

func ThereIsARepairOrderReceived(t *testing.T) uuid.UUID {
	t.Helper()
	return testsupport.ThereIsARepairOrderReceived(t, env.db)
}

// -----------------------------------------------------------------------------
// DB asserts (sem receber db)
// -----------------------------------------------------------------------------

func estimateItemsCountForRepairOrder(t *testing.T, repairOrderID uuid.UUID, expected int) {
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
