package product

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/tests/testsupport"

	productApp "github.com/soat13/fase-1-oficina/internal/product/application"
	productDB "github.com/soat13/fase-1-oficina/internal/product/infra/db"
	productHTTP "github.com/soat13/fase-1-oficina/internal/product/infra/http"
)

type httpTestApp struct {
	app *fiber.App
	tdb *testsupport.TestDB
}

// Sobe app + DB para CADA teste e fecha no Cleanup.
func setupHTTP(t *testing.T) *httpTestApp {
	t.Helper()

	tdb := testsupport.NewTestDB(t)

	// Repo + UCs
	repo := productDB.NewBunProductRepository(tdb.DB)
	createUC := productApp.NewCreateProduct(repo)
	updateUC := productApp.NewUpdateProduct(repo)
	deleteUC := productApp.NewDeleteProduct(repo)
	getUC := productApp.NewGetProduct(repo)
	listUC := productApp.NewListProducts(repo)

	// Fiber app
	app := fiber.New()
	app.Use(logger.New())

	// Handler + rotas (/admin/products/*)
	h := productHTTP.NewHandler(createUC, updateUC, deleteUC, getUC, listUC)
	productHTTP.Register(app, h)

	return &httpTestApp{app: app, tdb: tdb}
}

// Ambiente visível pelos testes
var env struct {
	app *fiber.App
	db  *bun.DB
}

// Inicializa app+db para CADA teste e garante fechamento via t.Cleanup
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
