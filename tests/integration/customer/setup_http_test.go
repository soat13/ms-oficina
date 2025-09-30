package customer

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/tests/testsupport"

	customerApp "github.com/soat13/fase-1-oficina/internal/customer/application"
	customerDB "github.com/soat13/fase-1-oficina/internal/customer/infra/db"
	customerHTTP "github.com/soat13/fase-1-oficina/internal/customer/infra/http"
)

type httpTestApp struct {
	app *fiber.App
	tdb *testsupport.TestDB
}

func setupHTTP(t *testing.T) *httpTestApp {
	t.Helper()

	tdb := testsupport.NewTestDB(t)

	repo := customerDB.NewBunCustomerRepository(tdb.DB)
	createUC := customerApp.NewCreateCustomer(repo)
	updateUC := customerApp.NewUpdateCustomer(repo)
	deleteUC := customerApp.NewDeleteCustomer(repo)
	getUC := customerApp.NewGetCustomer(repo)
	listUC := customerApp.NewListCustomers(repo)

	app := fiber.New()
	app.Use(logger.New())

	h := customerHTTP.NewHandler(createUC, updateUC, deleteUC, getUC, listUC)
	customerHTTP.Register(app, h)

	return &httpTestApp{app: app, tdb: tdb}
}

var env struct {
	app *fiber.App
	db  *bun.DB
}

func ensureSetup(t *testing.T) {
	t.Helper()
	ta := setupHTTP(t)
	env.app = ta.app
	env.db = ta.tdb.DB

	t.Cleanup(func() {
		ta.tdb.Close()
		env.app = nil
		env.db = nil
	})
}
