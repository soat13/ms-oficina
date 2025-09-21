package service

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/tests/testsupport"

	serviceApp "github.com/soat13/fase-1-oficina/internal/service/application"
	serviceDB "github.com/soat13/fase-1-oficina/internal/service/infra/db"
	serviceHTTP "github.com/soat13/fase-1-oficina/internal/service/infra/http"
)

type httpTestApp struct {
	app *fiber.App
	tdb *testsupport.TestDB
}

func setupHTTP(t *testing.T) *httpTestApp {
	t.Helper()

	tdb := testsupport.NewTestDB(t)

	repo := serviceDB.NewBunServiceRepository(tdb.DB)
	createUC := serviceApp.NewCreateService(repo)
	updateUC := serviceApp.NewUpdateService(repo)
	deleteUC := serviceApp.NewDeleteService(repo)
	getUC := serviceApp.NewGetService(repo)
	listUC := serviceApp.NewListServices(repo)

	app := fiber.New()
	app.Use(logger.New())

	h := serviceHTTP.NewHandler(createUC, updateUC, deleteUC, getUC, listUC)
	serviceHTTP.Register(app, h)

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
