package user

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/tests/testsupport"

	userApp "github.com/soat13/fase-1-oficina/internal/user/application"
	userDB "github.com/soat13/fase-1-oficina/internal/user/infra/db"
	userHTTP "github.com/soat13/fase-1-oficina/internal/user/infra/http"
	errorHelper "github.com/soat13/fase-1-oficina/pkg/error"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
)

type httpTestApp struct {
	app *fiber.App
	tdb *testsupport.TestDB
}

func setupHTTP(t *testing.T) *httpTestApp {
	t.Helper()

	tdb := testsupport.NewTestDB(t)

	repo := userDB.NewBunUserRepository(tdb.DB)
	createUC := userApp.NewCreateUser(repo)
	updateUC := userApp.NewUpdateUser(repo)
	deleteUC := userApp.NewDeleteUser(repo)
	getUC := userApp.NewGetUser(repo)
	listUC := userApp.NewListUsers(repo)

	app := fiber.New()
	app.Use(logger.New())

	errorResolver := errorHelper.NewErrorResolver()
	errorHandler := fiberHelper.NewErrorHandler(errorResolver)
	h := userHTTP.NewHandler(createUC, updateUC, deleteUC, getUC, listUC, errorHandler)
	userHTTP.Register(app, h)

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
