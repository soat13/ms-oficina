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

	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	estimateInfraDB "github.com/soat13/fase-1-oficina/internal/estimate/infra/db"
	estimateInfraHttp "github.com/soat13/fase-1-oficina/internal/estimate/infra/http"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application/listeners"
	repairOrderDB "github.com/soat13/fase-1-oficina/internal/repairorder/infra/db"
	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
	"github.com/soat13/fase-1-oficina/internal/shared/events/estimate"
)

func main() {
	loadEnv()

	db := newDB()
	defer db.bunDB.Close()

	repairOrderReader := estimateInfraDB.NewRepairOrderReader(db.bunDB)
	productCatalogReader := estimateInfraDB.NewProductCatalogReader(db.bunDB)
	serviceCatalogReader := estimateInfraDB.NewServiceCatalogReader(db.bunDB)
	estimateRepository := estimateInfraDB.NewBunEstimateRepository(db.bunDB)
	repairOrderRepository := repairOrderDB.NewBunRepairOrderRepository(db.bunDB)

	eventBus := eventbus.NewInMemoryBus()
	eventBus.Subscribe(estimate.Created{}.Topic(), listeners.OnEstimateCreated(repairOrderRepository))

	createEstimate := application.NewCreateEstimateFromRepairOrder(
		repairOrderReader,
		productCatalogReader,
		serviceCatalogReader,
		estimateRepository,
		eventBus,
	)

	app := newApp()
	estimateHttpHandler := estimateInfraHttp.NewHandler(createEstimate)
	estimateInfraHttp.Register(app, estimateHttpHandler)

	log.Println("✅ app iniciado na porta 8080")
	if err := app.Listen(":8080"); err != nil {
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
	return app
}
