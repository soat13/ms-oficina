package testsupport

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/document"
	passwordVO "github.com/soat13/fase-1-oficina/pkg/valueobjects/password"
)

var (
	OilFilterID       = uuid.MustParse("1f09479b-6074-69d0-a0d0-9a7c61ccc8bb")
	EngineOilChangeID = uuid.MustParse("1f094799-f210-63a0-a800-e67186b3a9f6")
)

func ThereIsAFinishedRepairOrder(t *testing.T, db *bun.DB) uuid.UUID {
	t.Helper()
	return ThereIsARepairOrderWithStatus(t, db, repairorder.StatusFinished)
}

func ThereIsAnApprovedRepairOrder(t *testing.T, db *bun.DB) uuid.UUID {
	t.Helper()
	return ThereIsARepairOrderWithStatus(t, db, repairorder.StatusApproved)
}

func ThereIsARepairOrderInDiagnostics(t *testing.T, db *bun.DB) uuid.UUID {
	t.Helper()
	return ThereIsARepairOrderWithStatus(t, db, repairorder.StatusInDiagnostics)
}

func ThereIsARepairOrderInExecution(t *testing.T, db *bun.DB) uuid.UUID {
	t.Helper()
	return ThereIsARepairOrderWithStatus(t, db, repairorder.StatusInExecution)
}

func ThereIsAReceivedRepairOrder(t *testing.T, db *bun.DB) uuid.UUID {
	t.Helper()
	return ThereIsARepairOrderWithStatus(t, db, repairorder.StatusReceived)
}

func ThereIsARepairOrderWithStatus(t *testing.T, db *bun.DB, status repairorder.Status) uuid.UUID {
	t.Helper()
	tag := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))[:8]
	name := "Customer " + tag
	plate := "T" + tag[:6]
	doc := GenerateDocument()
	email := faker.Email()
	phoneNumber := "11987654321"

	customerID := ThereIsACustomerWithID(t, db, uuid.Nil, name, doc, phoneNumber, email)
	vehicleID := ThereIsAVehicle(t, db, uuid.Nil, customerID, plate, "Toyota", "Corolla", 2020)
	repairOrderID := ThereIsARepairOrder(t, db, uuid.Nil, customerID, vehicleID, string(status))

	return repairOrderID
}

func ThereIsACustomerWithID(t *testing.T, db *bun.DB, id uuid.UUID, name string, doc document.Document, phoneNumber, email string) uuid.UUID {
	t.Helper()
	if id == uuid.Nil {
		id = uuid.New()
	}

	_, err := db.NewRaw(`
		INSERT INTO customers (id, name, document, document_type, phone_number, email)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING
	`, id, name, doc.Value, doc.Type(), phoneNumber, email).Exec(context.Background())
	require.NoError(t, err, "falha ao inserir customer")
	return id
}

func ThereIsACustomerWithDocument(t *testing.T, db *bun.DB, id uuid.UUID, name, documentStr, documentType, phoneNumber, email string) uuid.UUID {
	t.Helper()
	if id == uuid.Nil {
		id = uuid.New()
	}

	doc, err := document.New(documentStr)
	require.NoError(t, err, "invalid document in test")
	normalizedDoc := doc.Value

	_, err = db.NewRaw(`
		INSERT INTO customers (id, name, document, document_type, phone_number, email)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING
	`, id, name, normalizedDoc, documentType, phoneNumber, email).Exec(context.Background())
	require.NoError(t, err, "falha ao inserir customer")
	return id
}

func ThereIsACustomer(t *testing.T, db *bun.DB, id uuid.UUID, name, phoneNumber, email string) uuid.UUID {
	t.Helper()
	if id == uuid.Nil {
		id = uuid.New()
	}

	doc := GenerateDocument()

	_, err := db.NewRaw(`
		INSERT INTO customers (id, name, document, document_type, phone_number, email)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO NOTHING
	`, id, name, doc.Value, doc.Type(), phoneNumber, email).Exec(context.Background())
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

func ThereIsAnEstimateForRepairOrder(t *testing.T, container *bootstrap.Container, repairOrderID uuid.UUID) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	estimateID := uuid.New()
	_, err := container.DB.NewRaw(`
		INSERT INTO estimates (id, repair_order_id, status)
		VALUES (?, ?, ?)
	`, estimateID, repairOrderID, "awaiting_approval").Exec(ctx)
	require.NoError(t, err, "falha ao inserir estimate para o RO")

	_, err = container.DB.NewRaw(`
		INSERT INTO estimate_items (id, estimate_id, item_id, item_name, item_type, price, quantity)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, uuid.New(), estimateID, uuid.New(), "Oil Filter", "product", 1000, 1).Exec(ctx)
	require.NoError(t, err, "falha ao inserir estimate_item")

	return estimateID
}

func ThereIsAUser(t *testing.T, db *bun.DB, id uuid.UUID, name, documentStr, phoneNumber, email, password string, roles []string) uuid.UUID {
	t.Helper()
	if id == uuid.Nil {
		id = uuid.New()
	}

	doc, err := document.New(documentStr)
	require.NoError(t, err, "invalid document in test")
	normalizedDoc := doc.Value

	hashedPassword := hashPassword(t, password)
	_, err = db.NewRaw(`
		INSERT INTO users (id, name, document, document_type, email, phone_number, password, roles)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (document) DO NOTHING
	`, id, name, normalizedDoc, doc.Type(), email, phoneNumber, hashedPassword, pq.Array(roles)).Exec(context.Background())
	require.NoError(t, err, "failed to insert user")
	return id
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	vo, err := passwordVO.New(password)
	require.NoError(t, err, "invalid password in test")
	return vo.Hash
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

func GenerateDocument() document.Document {
	var documentValue string

	if rand.Intn(2) == 0 {
		documentValue = GenerateCPF()
	} else {
		documentValue = GenerateCNPJ()
	}

	doc, _ := document.New(documentValue)
	return doc
}

func GenerateCPF() string {
	var cpf [9]int
	for i := 0; i < 9; i++ {
		cpf[i] = rand.Intn(10)
	}

	// first digit
	sum := 0
	for i, j := 0, 10; i < 9; i, j = i+1, j-1 {
		sum += cpf[i] * j
	}
	d1 := (sum * 10) % 11
	if d1 == 10 {
		d1 = 0
	}

	// second digit
	sum = 0
	for i, j := 0, 11; i < 9; i, j = i+1, j-1 {
		sum += cpf[i] * j
	}
	sum += d1 * 2
	d2 := (sum * 10) % 11
	if d2 == 10 {
		d2 = 0
	}

	return fmt.Sprintf("%d%d%d.%d%d%d.%d%d%d-%d%d",
		cpf[0], cpf[1], cpf[2],
		cpf[3], cpf[4], cpf[5],
		cpf[6], cpf[7], cpf[8],
		d1, d2,
	)
}

func GenerateCNPJ() string {
	var cnpj [12]int
	for i := 0; i < 12; i++ {
		cnpj[i] = rand.Intn(10)
	}

	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	// first digit
	sum := 0
	for i := 0; i < 12; i++ {
		sum += cnpj[i] * weights1[i]
	}
	d1 := sum % 11
	if d1 < 2 {
		d1 = 0
	} else {
		d1 = 11 - d1
	}

	// secondo digit
	sum = 0
	for i := 0; i < 12; i++ {
		sum += cnpj[i] * weights2[i]
	}
	sum += d1 * weights2[12]
	d2 := sum % 11
	if d2 < 2 {
		d2 = 0
	} else {
		d2 = 11 - d2
	}

	return fmt.Sprintf("%d%d.%d%d%d.%d%d%d/%d%d%d%d-%d%d",
		cnpj[0], cnpj[1],
		cnpj[2], cnpj[3], cnpj[4],
		cnpj[5], cnpj[6], cnpj[7],
		cnpj[8], cnpj[9], cnpj[10], cnpj[11],
		d1, d2,
	)
}
