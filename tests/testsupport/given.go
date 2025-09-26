package testsupport

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

var (
	OilFilterID       = uuid.MustParse("1f09479b-6074-69d0-a0d0-9a7c61ccc8bb")
	EngineOilChangeID = uuid.MustParse("1f094799-f210-63a0-a800-e67186b3a9f6")
)

func ThereIsARepairOrderInDiagnostics(t *testing.T, db *bun.DB) uuid.UUID {
	t.Helper()
	return ThereIsARepairOrderWithStatus(t, db, repairorder.StatusInDiagnosis)
}

func ThereIsARepairOrderReceived(t *testing.T, db *bun.DB) uuid.UUID {
	t.Helper()
	return ThereIsARepairOrderWithStatus(t, db, repairorder.StatusReceived)
}

func ThereIsARepairOrderWithStatus(t *testing.T, db *bun.DB, status repairorder.Status) uuid.UUID {
	t.Helper()

	// todo: implement faker for names, cpfs, plates, etc.
	tag := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))[:8]
	cpf := "CPF" + tag
	name := "Customer " + tag
	plate := "T" + tag[:6]
	documentType := "CPF"
	cellphone := "11987654321"
	email := "customer@example.com"

	customerID := ThereIsACustomerWithID(t, db, uuid.Nil, name, cpf, documentType, cellphone, email)
	vehicleID := ThereIsAVehicle(t, db, uuid.Nil, customerID, plate, "Toyota", "Corolla", 2020)
	repairOrderID := ThereIsARepairOrder(t, db, uuid.Nil, customerID, vehicleID, string(status))

	return repairOrderID
}

func ThereIsACustomerWithID(t *testing.T, db *bun.DB, id uuid.UUID, name, document, documentType, cellphone, email string) uuid.UUID {
	t.Helper()
	if id == uuid.Nil {
		id = uuid.New()
	}
	_, err := db.NewRaw(`
		INSERT INTO customers (id, name, document, document_type, cellphone, email)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING
	`, id, name, document, documentType, cellphone, email).Exec(context.Background())
	require.NoError(t, err, "falha ao inserir customer")
	return id
}

func ThereIsACustomerWithDocument(t *testing.T, db *bun.DB, id uuid.UUID, name, document, documentType, cellphone, email string) uuid.UUID {
	t.Helper()
	if id == uuid.Nil {
		id = uuid.New()
	}
	_, err := db.NewRaw(`
		INSERT INTO customers (id, name, document, document_type, cellphone, email)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING
	`, id, name, document, documentType, cellphone, email).Exec(context.Background())
	require.NoError(t, err, "falha ao inserir customer")
	return id
}

func ThereIsAVehicle(t *testing.T, db *bun.DB, id, customerID uuid.UUID, plate, brand, model string, year int) uuid.UUID {
	t.Helper()
	if id == uuid.Nil {
		id = uuid.New()
	}
	_, err := db.NewRaw(`
		INSERT INTO vehicles (id, customer_id, plate, brand, model, year)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING
	`, id, customerID, plate, brand, model, year).Exec(context.Background())
	require.NoError(t, err, "falha ao inserir vehicle")
	return id
}

func ThereIsAService(t *testing.T, db *bun.DB, id uuid.UUID, name string, price int64) uuid.UUID {
	t.Helper()
	if id == uuid.Nil {
		id = uuid.New()
	}
	require.NoError(t, insertService(context.Background(), db, id, name, price))
	return id
}

func ThereIsAProduct(t *testing.T, db *bun.DB, id uuid.UUID, name string, priceCents int64, stock int) uuid.UUID {
	t.Helper()
	if id == uuid.Nil {
		id = uuid.New()
	}
	require.NoError(t, insertProduct(context.Background(), db, id, name, priceCents, stock))
	return id
}

func ThereIsARepairOrder(t *testing.T, db *bun.DB, id, customerID, vehicleID uuid.UUID, status string) uuid.UUID {
	t.Helper()
	if id == uuid.Nil {
		id = uuid.New()
	}
	_, err := db.NewRaw(`
		INSERT INTO repair_orders (id, customer_id, vehicle_id, status)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING
	`, id, customerID, vehicleID, status).Exec(context.Background())
	require.NoError(t, err, "falha ao inserir repair_order")
	return id
}

func insertService(ctx context.Context, db *bun.DB, id uuid.UUID, name string, cents int64) error {

	_, err := db.NewRaw(`
		INSERT INTO services (id, name, price)
		VALUES (?, ?, ?)
		ON CONFLICT (id) DO NOTHING
	`, id, name, cents).Exec(ctx)
	return err
}

func insertProduct(ctx context.Context, db *bun.DB, id uuid.UUID, name string, priceCents int64, stock int) error {
	_, err := db.NewRaw(`
		INSERT INTO products (id, name, price, stock)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING
	`, id, name, priceCents, stock).Exec(ctx)
	return err
}
