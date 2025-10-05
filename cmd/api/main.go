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
	estimaterListeners "github.com/soat13/fase-1-oficina/internal/estimate/application/listeners"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
	productEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/product"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
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
	serviceDocs "github.com/soat13/fase-1-oficina/internal/service/infra/docs"
	serviceHTTP "github.com/soat13/fase-1-oficina/internal/service/infra/http"

	// customer
	customerApp "github.com/soat13/fase-1-oficina/internal/customer/application"
	customerDB "github.com/soat13/fase-1-oficina/internal/customer/infra/db"
	customerHTTP "github.com/soat13/fase-1-oficina/internal/customer/infra/http"

	// user
	userApp "github.com/soat13/fase-1-oficina/internal/user/application"
	userDB "github.com/soat13/fase-1-oficina/internal/user/infra/db"
	userHTTP "github.com/soat13/fase-1-oficina/internal/user/infra/http"

	// repair order listeners
	repairOrderApp "github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application/listeners"
	repairOrderDB "github.com/soat13/fase-1-oficina/internal/repairorder/infra/db"
	repairOrderHTTP "github.com/soat13/fase-1-oficina/internal/repairorder/infra/http"

	// shared
	"github.com/soat13/fase-1-oficina/internal/shared/errors"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
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
	errorResolver.RegisterHTTPBadRequestError(errors.ErrInvalidJSON)
	errorResolver.RegisterHTTPConflictError(errors.ErrInvalidStatusTransaction)
	errorResolver.RegisterHTTPNotFoundError(sharedRepairOrder.ErrRepairOrderNotFound)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrInvalidDocument)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrInvalidPhoneNumber)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrInvalidEmail)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrPasswordTooShort)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrPasswordTooLong)
	errorResolver.RegisterHTTPUnprocessableError(errors.ErrInvalidPasswordHash)

	fiberApp := newApp()
	serviceDocs.Register(fiberApp)
	errorHandler := fiberHelper.NewErrorHandler(errorResolver)

	// -----------------------------------------------------------------------------
	// Estimate wiring
	// -----------------------------------------------------------------------------
	repairOrderReader := estimateInfraDB.NewRepairOrderReader(db.bunDB)
	productCatalogReader := estimateInfraDB.NewProductCatalogReader(db.bunDB)
	serviceCatalogReader := estimateInfraDB.NewServiceCatalogReader(db.bunDB)
	estimateRepository := estimateInfraDB.NewBunRepository(db.bunDB)

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

	eventBus.Subscribe(
		productEvent.StockReduceConfirmed{}.Topic(),
		estimaterListeners.OnStockReduceConfirmed(estimateRepository, eventBus),
	)

	// -----------------------------------------------------------------------------
	// Services wiring
	// -----------------------------------------------------------------------------
	serviceRepository := serviceDB.NewBunRepository(db.bunDB)
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
	releaseVehicle := repairOrderApp.NewReleaseVehicle(repairOrderRepository)

	repairOrderHTTPHandler := repairOrderHTTP.NewHandler(startExecution, finishExecution, releaseVehicle, errorHandler)
	repairOrderHTTP.Register(fiberApp, repairOrderHTTPHandler)

	// Event listeners
	eventBus.Subscribe(estimate.Created{}.Topic(), listeners.OnEstimateCreated(repairOrderRepository))
	eventBus.Subscribe(estimate.Approved{}.Topic(), listeners.OnEstimateApproved(repairOrderRepository))

	// -----------------------------------------------------------------------------
	// Customers wiring
	// -----------------------------------------------------------------------------
	customerHandler := customerHTTP.NewHandler(createCus, updateCus, deleteCus, getCus, listCus, errorHandler)
	customerHTTP.Register(fiberApp, customerHandler)

	// -----------------------------------------------------------------------------
	// User wiring
	// -----------------------------------------------------------------------------
	userRepo := userDB.NewBunUserRepository(db.bunDB)
	createUser := userApp.NewCreateUser(userRepo)
	updateUser := userApp.NewUpdateUser(userRepo)
	deleteUser := userApp.NewDeleteUser(userRepo)
	getUser := userApp.NewGetUser(userRepo)
	listUser := userApp.NewListUsers(userRepo)

	userHandler := userHTTP.NewHandler(createUser, updateUser, deleteUser, getUser, listUser, errorHandler)
	userHTTP.Register(fiberApp, userHandler)

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
	// TODO: adicionar middleware de JWT e aplicar no grupo /admin/*
	return app
}
