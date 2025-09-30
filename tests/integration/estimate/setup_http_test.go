package estimate

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/uptrace/bun"

	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	estimateInfraDB "github.com/soat13/fase-1-oficina/internal/estimate/infra/db"
	estimateInfraHTTP "github.com/soat13/fase-1-oficina/internal/estimate/infra/http"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application/listeners"
	repairOrderDB "github.com/soat13/fase-1-oficina/internal/repairorder/infra/db"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	estimateev "github.com/soat13/fase-1-oficina/internal/shared/events/estimate"
	"github.com/soat13/fase-1-oficina/tests/testsupport"
)

type httpTestApp struct {
	app *fiber.App
	tdb *testsupport.TestDB
}

func setupHTTP(t *testing.T) *httpTestApp {
	t.Helper()

	tdb := testsupport.NewTestDB(t)

	repairOrderReader := estimateInfraDB.NewRepairOrderReader(tdb.DB)
	productCatalogReader := estimateInfraDB.NewProductCatalogReader(tdb.DB)
	serviceCatalogReader := estimateInfraDB.NewServiceCatalogReader(tdb.DB)
	estimateRepository := estimateInfraDB.NewBunEstimateRepository(tdb.DB)
	repairOrderRepository := repairOrderDB.NewBunRepairOrderRepository(tdb.DB)

	bus := eventbus.NewInMemoryBus()
	bus.Subscribe(estimateev.Created{}.Topic(), listeners.OnEstimateCreated(repairOrderRepository))
	bus.Subscribe(estimateev.ApprovedByCustomer{}.Topic(), listeners.OnEstimateApprovedByCustomer(repairOrderRepository))

	createEstimate := application.NewCreateEstimate(
		repairOrderReader,
		productCatalogReader,
		serviceCatalogReader,
		estimateRepository,
		bus,
	)

	approveEstimate := application.NewApproveEstimate(estimateRepository, bus)

	app := fiber.New()
	app.Use(logger.New())

	h := estimateInfraHTTP.NewHandler(createEstimate, approveEstimate)
	estimateInfraHTTP.Register(app, h)

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
