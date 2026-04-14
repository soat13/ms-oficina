package repairorder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/tests/testsupport"
	testauth "github.com/soat13/fase-1-oficina/tests/testsupport/auth"
)

// -----------------------------------------------------------------------------
// Given
// -----------------------------------------------------------------------------

func givenProductExists(ctx context.Context, nome string, estoque int) error {
	w := worldFromContext(ctx)
	id := uuid.New()
	_, err := w.setup.Container.DB.NewRaw(`
		INSERT INTO products (id, name, price, stock)
		VALUES (?, ?, ?, ?)
	`, id, strings.ToLower(nome), int64(10000), estoque).Exec(ctx)
	if err != nil {
		return fmt.Errorf("seed product %q: %w", nome, err)
	}
	w.productIDs[nome] = id
	return nil
}

func givenServiceExists(ctx context.Context, nome string) error {
	w := worldFromContext(ctx)
	id := uuid.New()
	_, err := w.setup.Container.DB.NewRaw(`
		INSERT INTO services (id, name, price)
		VALUES (?, ?, ?)
	`, id, strings.ToLower(nome), int64(10000)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("seed service %q: %w", nome, err)
	}
	w.serviceIDs[nome] = id
	return nil
}

func givenRepairOrderExistsWithStatus(ctx context.Context, status string) error {
	w := worldFromContext(ctx)
	customerID, err := insertCustomer(ctx, w)
	if err != nil {
		return err
	}
	vehicleID, err := insertVehicle(ctx, w, customerID)
	if err != nil {
		return err
	}
	roID := uuid.New()
	_, err = w.setup.Container.DB.NewRaw(`
		INSERT INTO repair_orders (id, customer_id, vehicle_id, status)
		VALUES (?, ?, ?, ?)
	`, roID, customerID, vehicleID, status).Exec(ctx)
	if err != nil {
		return fmt.Errorf("seed repair_order (%s): %w", status, err)
	}
	w.repairOrderID = roID
	w.initialStatus = status
	return nil
}

func insertCustomer(ctx context.Context, w *world) (uuid.UUID, error) {
	id := uuid.New()
	tag := strings.ReplaceAll(id.String(), "-", "")[:8]
	cpf := testsupport.GenerateCPF()
	_, err := w.setup.Container.DB.NewRaw(`
		INSERT INTO customers (id, name, document, document_type, phone_number, email)
		VALUES (?, ?, ?, ?, ?, ?)
	`, id, "Customer "+tag, cpf, "CPF", "11987654321", tag+"@example.com").Exec(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("seed customer: %w", err)
	}
	return id, nil
}

func insertVehicle(ctx context.Context, w *world, customerID uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	tag := strings.ReplaceAll(id.String(), "-", "")[:6]
	_, err := w.setup.Container.DB.NewRaw(`
		INSERT INTO vehicles (id, customer_id, plate, brand, model, year)
		VALUES (?, ?, ?, ?, ?, ?)
	`, id, customerID, strings.ToUpper("T"+tag), "Toyota", "Corolla", 2020).Exec(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("seed vehicle: %w", err)
	}
	return id, nil
}

// -----------------------------------------------------------------------------
// When
// -----------------------------------------------------------------------------

func whenStartDiagnostics(ctx context.Context) error {
	w := worldFromContext(ctx)
	return callEndpoint(ctx, w, http.MethodPost,
		"/admin/repair-orders/"+w.repairOrderID.String()+"/start-diagnostics", nil)
}

func whenFinishDiagnosticsWith(ctx context.Context, table *godog.Table) error {
	w := worldFromContext(ctx)

	products := []map[string]any{}
	services := []map[string]any{}

	for _, row := range table.Rows[1:] {
		tipo := strings.ToLower(strings.TrimSpace(row.Cells[0].Value))
		nome := strings.TrimSpace(row.Cells[1].Value)
		qtd, err := strconv.Atoi(strings.TrimSpace(row.Cells[2].Value))
		if err != nil {
			return fmt.Errorf("invalid quantity %q: %w", row.Cells[2].Value, err)
		}

		switch tipo {
		case "produto":
			id, ok := w.productIDs[nome]
			if !ok {
				return fmt.Errorf("product %q was not seeded in Background", nome)
			}
			products = append(products, map[string]any{"id": id.String(), "quantity": qtd})
		case "servico", "serviço":
			id, ok := w.serviceIDs[nome]
			if !ok {
				return fmt.Errorf("service %q was not seeded in Background", nome)
			}
			services = append(services, map[string]any{"id": id.String(), "quantity": qtd})
		default:
			return fmt.Errorf("unknown type %q (use 'produto' or 'servico')", tipo)
		}
	}

	body := map[string]any{"products": products, "services": services}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal body: %w", err)
	}

	return callEndpoint(ctx, w, http.MethodPost,
		"/admin/repair-orders/"+w.repairOrderID.String()+"/finish-diagnostics", payload)
}

// -----------------------------------------------------------------------------
// Then
// -----------------------------------------------------------------------------

func thenHttpResponseShouldBe(ctx context.Context, expected int) error {
	w := worldFromContext(ctx)
	if w.lastResp == nil {
		return fmt.Errorf("no HTTP call has been made yet")
	}
	if w.lastResp.StatusCode != expected {
		return fmt.Errorf("expected HTTP %d, got %d", expected, w.lastResp.StatusCode)
	}
	return nil
}

func thenRepairOrderShouldBeInStatus(ctx context.Context, expected string) error {
	w := worldFromContext(ctx)
	var status string
	err := w.setup.Container.DB.NewRaw(
		`SELECT status FROM repair_orders WHERE id = ?`, w.repairOrderID,
	).Scan(ctx, &status)
	if err != nil {
		return fmt.Errorf("read repair_order status: %w", err)
	}
	if status != expected {
		return fmt.Errorf("expected status %q, got %q", expected, status)
	}
	return nil
}

func thenEstimateShouldBeInStatus(ctx context.Context, expected string) error {
	w := worldFromContext(ctx)
	var status string
	err := w.setup.Container.DB.NewRaw(
		`SELECT status FROM estimates WHERE repair_order_id = ?`, w.repairOrderID,
	).Scan(ctx, &status)
	if err != nil {
		return fmt.Errorf("no estimate found for repair_order: %w", err)
	}
	if status != expected {
		return fmt.Errorf("expected estimate status %q, got %q", expected, status)
	}
	return nil
}

func thenEstimateShouldContainNItems(ctx context.Context, expected int) error {
	w := worldFromContext(ctx)
	var count int
	err := w.setup.Container.DB.NewRaw(`
		SELECT COUNT(*) FROM estimate_items ei
		JOIN estimates e ON e.id = ei.estimate_id
		WHERE e.repair_order_id = ?
	`, w.repairOrderID).Scan(ctx, &count)
	if err != nil {
		return fmt.Errorf("count estimate items: %w", err)
	}
	if count != expected {
		return fmt.Errorf("expected %d items, got %d", expected, count)
	}
	return nil
}

// -----------------------------------------------------------------------------
// HTTP helper
// -----------------------------------------------------------------------------

func callEndpoint(_ context.Context, w *world, method, path string, body []byte) error {
	var reader *bytes.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	testauth.AddAuthHeader(req, w.setup.AuthToken)

	resp, err := w.setup.Container.FiberApp.Test(req, -1)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	w.lastResp = resp
	return nil
}
