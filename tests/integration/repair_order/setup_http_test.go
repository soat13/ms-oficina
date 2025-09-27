package repairorder

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/soat13/fase-1-oficina/internal/shared/errors"
	errorHelper "github.com/soat13/fase-1-oficina/pkg/error"
	"github.com/uptrace/bun"

	roapp "github.com/soat13/fase-1-oficina/internal/repairorder/application"
	repairOrderDB "github.com/soat13/fase-1-oficina/internal/repairorder/infra/db"
	repairOrderHTTP "github.com/soat13/fase-1-oficina/internal/repairorder/infra/http"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"

	"github.com/soat13/fase-1-oficina/tests/testsupport"
)

type httpTestApp struct {
	app *fiber.App
	tdb *testsupport.TestDB
}

func setupHTTP(t *testing.T) *httpTestApp {
	t.Helper()

	tdb := testsupport.NewTestDB(t)

	roRepo := repairOrderDB.NewBunRepairOrderRepository(tdb.DB)
	startExecUC := roapp.NewStartExecution(roRepo)

	app := fiber.New()
	app.Use(logger.New())

	errorResolver := errorHelper.NewErrorResolver()

	errorResolver.RegisterHTTPBadRequestError(errors.ErrInvalidID)
	errorResolver.RegisterHTTPConflictError(errors.ErrInvalidStatusTransaction)

	errorHandler := fiberHelper.NewErrorHandler(errorResolver)
	h := repairOrderHTTP.NewHandler(startExecUC, errorHandler)
	repairOrderHTTP.Register(app, h)

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
