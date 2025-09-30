package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	// pkg
	errorHelper "github.com/soat13/fase-1-oficina/pkg/error"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"

	// estimate
	estimateApp "github.com/soat13/fase-1-oficina/internal/estimate/application"
	estimateInfraDB "github.com/soat13/fase-1-oficina/internal/estimate/infra/db"
	estimateInfraHttp "github.com/soat13/fase-1-oficina/internal/estimate/infra/http"

	// services
	serviceApp "github.com/soat13/fase-1-oficina/internal/service/application"
	serviceDB "github.com/soat13/fase-1-oficina/internal/service/infra/db"
	serviceHTTP "github.com/soat13/fase-1-oficina/internal/service/infra/http"
    serviceDocs "github.com/soat13/fase-1-oficina/internal/service/infra/docs"

	// customer
	customerApp "github.com/soat13/fase-1-oficina/internal/customer/application"
	customerDB "github.com/soat13/fase-1-oficina/internal/customer/infra/db"
	customerHTTP "github.com/soat13/fase-1-oficina/internal/customer/infra/http"

	// repair order listeners
	repairOrderApp "github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application/listeners"
	repairOrderDB "github.com/soat13/fase-1-oficina/internal/repairorder/infra/db"
	repairOrderHTTP "github.com/soat13/fase-1-oficina/internal/repairorder/infra/http"

	// shared
	"github.com/soat13/fase-1-oficina/internal/shared/errors"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	"github.com/soat13/fase-1-oficina/internal/shared/events/estimate"
)

func main() {
	loadEnv()

	db := newDB()
	defer db.bunDB.Close()

	eventBus := eventbus.NewInMemoryBus()

	// -----------------------------------------------------------------------------
	// HTTP server setup
	// -----------------------------------------------------------------------------
	errorResolver := errorHelper.NewErrorResolver()

	errorResolver.RegisterHTTPBadRequestError(errors.ErrInvalidID)
	errorResolver.RegisterHTTPConflictError(errors.ErrInvalidStatusTransaction)

	// -----------------------------------------------------------------------------
	// HTTP server setup
	// -----------------------------------------------------------------------------
	fiberApp := newApp()
	errorHandler := fiberHelper.NewErrorHandler(errorResolver)

	errorResolver.RegisterHTTPBadRequestError(errors.ErrInvalidID)
	errorResolver.RegisterHTTPConflictError(errors.ErrInvalidStatusTransaction)

	// -----------------------------------------------------------------------------
	// Estimate wiring
	// -----------------------------------------------------------------------------
	repairOrderReader := estimateInfraDB.NewRepairOrderReader(db.bunDB)
	productCatalogReader := estimateInfraDB.NewProductCatalogReader(db.bunDB)
	serviceCatalogReader := estimateInfraDB.NewServiceCatalogReader(db.bunDB)
	estimateRepository := estimateInfraDB.NewBunEstimateRepository(db.bunDB)

	createEstimate := estimateApp.NewCreateEstimate(
		repairOrderReader,
		productCatalogReader,
		serviceCatalogReader,
		estimateRepository,
		eventBus,
	)

	approveEstimate := estimateApp.NewApproveEstimate(estimateRepository, eventBus)

	estimateHttpHandler := estimateInfraHttp.NewHandler(createEstimate, approveEstimate)
	estimateInfraHttp.Register(fiberApp, estimateHttpHandler)

	// -----------------------------------------------------------------------------
	// Services wiring
	// -----------------------------------------------------------------------------
	serviceRepository := serviceDB.NewBunServiceRepository(db.bunDB)
	createService := serviceApp.NewCreateService(serviceRepository)
	updateService := serviceApp.NewUpdateService(serviceRepository)
	deleteService := serviceApp.NewDeleteService(serviceRepository)
	getService := serviceApp.NewGetService(serviceRepository)
	listService := serviceApp.NewListServices(serviceRepository)

	serviceHttpHandler := serviceHTTP.NewHandler(createService, updateService, deleteService, getService, listService)
	serviceHTTP.Register(fiberApp, serviceHttpHandler)

	// -----------------------------------------------------------------------------
	// Customer wiring
	// -----------------------------------------------------------------------------
	customerRepo := customerDB.NewBunCustomerRepository(db.bunDB)
	createCus := customerApp.NewCreateCustomer(customerRepo)
	updateCus := customerApp.NewUpdateCustomer(customerRepo)
	deleteCus := customerApp.NewDeleteCustomer(customerRepo)
	getCus := customerApp.NewGetCustomer(customerRepo)
	listCus := customerApp.NewListCustomers(customerRepo)

	// -----------------------------------------------------------------------------
	// Repairorder wiring
	// -----------------------------------------------------------------------------
	repairOrderRepository := repairOrderDB.NewBunRepairOrderRepository(db.bunDB)
	startExecution := repairOrderApp.NewStartExecution(repairOrderRepository)
	finishExecution := repairOrderApp.NewFinishExecution(repairOrderRepository)

	repairOrderHTTPHandler := repairOrderHTTP.NewHandler(startExecution, finishExecution, errorHandler)
	repairOrderHTTP.Register(fiberApp, repairOrderHTTPHandler)

	eventBus.Subscribe(estimate.Created{}.Topic(), listeners.OnEstimateCreated(repairOrderRepository))
	eventBus.Subscribe(estimate.ApprovedByCustomer{}.Topic(), listeners.OnEstimateApprovedByCustomer(repairOrderRepository))
	// todo: add approvedByCustomer product subscriber to reduce stock

	//customers
	customerHandler := customerHTTP.NewHandler(createCus, updateCus, deleteCus, getCus, listCus)
	customerHTTP.Register(fiberApp, customerHandler)

	// -----------------------------------------------------------------------------
	// HTTP server start
	// -----------------------------------------------------------------------------
	port := os.Getenv("PORT")
	if err := fiberApp.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env não encontrado, usando variáveis de ambiente do sistema")
	}
}

type DB struct {
	SQL   *sql.DB
	bunDB *bun.DB
}

func newDB() *DB {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		panic("env PG_DSN not found")
	}

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := sqlDB.PingContext(context.Background()); err != nil {
		log.Fatal(err)
	}

	bunDB := bun.NewDB(sqlDB, pgdialect.New())
	return &DB{SQL: sqlDB, bunDB: bunDB}
}

func newApp() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
  serviceDocs.Register(app)
	// TODO: adicionar middleware de JWT e aplicar no grupo /admin/*
	return app
}
