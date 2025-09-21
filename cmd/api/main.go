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

	// estimate
	estimateApp "github.com/soat13/fase-1-oficina/internal/estimate/application"
	estimateInfraDB "github.com/soat13/fase-1-oficina/internal/estimate/infra/db"
	estimateInfraHttp "github.com/soat13/fase-1-oficina/internal/estimate/infra/http"

	// services
	serviceApp "github.com/soat13/fase-1-oficina/internal/service/application"
	serviceDB "github.com/soat13/fase-1-oficina/internal/service/infra/db"
	serviceHTTP "github.com/soat13/fase-1-oficina/internal/service/infra/http"

	// customer

	// shared
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	"github.com/soat13/fase-1-oficina/internal/shared/events/estimate"

	// repair order listeners
	"github.com/soat13/fase-1-oficina/internal/repairorder/application/listeners"
	repairOrderDB "github.com/soat13/fase-1-oficina/internal/repairorder/infra/db"
)

func main() {
	loadEnv()

	db := newDB()
	defer db.bunDB.Close()

	eventBus := eventbus.NewInMemoryBus()

	// -----------------------------------------------------------------------------
	// Estimate wiring
	// -----------------------------------------------------------------------------
	repairOrderReader := estimateInfraDB.NewRepairOrderReader(db.bunDB)
	productCatalogReader := estimateInfraDB.NewProductCatalogReader(db.bunDB)
	serviceCatalogReader := estimateInfraDB.NewServiceCatalogReader(db.bunDB)
	estimateRepository := estimateInfraDB.NewBunEstimateRepository(db.bunDB)
	repairOrderRepository := repairOrderDB.NewBunRepairOrderRepository(db.bunDB)

	eventBus.Subscribe(estimate.Created{}.Topic(), listeners.OnEstimateCreated(repairOrderRepository))
	eventBus.Subscribe(estimate.Approved{}.Topic(), listeners.OnEstimateApproved(repairOrderRepository))

	createEstimate := estimateApp.NewCreateEstimateFromRepairOrder(
		repairOrderReader,
		productCatalogReader,
		serviceCatalogReader,
		estimateRepository,
		eventBus,
	)

	// -----------------------------------------------------------------------------
	// Services wiring
	// -----------------------------------------------------------------------------
	serviceRepo := serviceDB.NewBunServiceRepository(db.bunDB)
	createSvc := serviceApp.NewCreateService(serviceRepo)
	updateSvc := serviceApp.NewUpdateService(serviceRepo)
	deleteSvc := serviceApp.NewDeleteService(serviceRepo)
	getSvc := serviceApp.NewGetService(serviceRepo)
	listSvc := serviceApp.NewListServices(serviceRepo)

	// TODO: add customer wiring
	// -----------------------------------------------------------------------------
	// Customer wiring
	// -----------------------------------------------------------------------------
	// customerRepo := customerDB.NewBunCustomerRepository(db.bunDB)
	// createCus := customerApp.NewCreateCustomer(customerRepo)
	// updateCus := customerApp.NewUpdateCustomer(customerRepo)
	// deleteCus := customerApp.NewDeleteCustomer(customerRepo)
	// getCus := customerApp.NewGetCustomer(customerRepo)
	// listCus := customerApp.NewListCustomers(customerRepo)

	// -----------------------------------------------------------------------------
	// HTTP app & routes
	// -----------------------------------------------------------------------------
	app := newApp()

	// estimate
	estimateHttpHandler := estimateInfraHttp.NewHandler(createEstimate)
	estimateInfraHttp.Register(app, estimateHttpHandler)

	// services
	serviceHttpHandler := serviceHTTP.NewHandler(createSvc, updateSvc, deleteSvc, getSvc, listSvc)
	serviceHTTP.Register(app, serviceHttpHandler)

	// -----------------------------------------------------------------------------
	// HTTP server start
	// -----------------------------------------------------------------------------
	port := os.Getenv("PORT")
	log.Println("✅ app iniciado na porta: " + port)
	if err := app.Listen(":" + port); err != nil {
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
