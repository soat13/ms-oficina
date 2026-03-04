package bootstrap

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/stdlib"
	infraMessaging "github.com/soat13/fase-1-oficina/internal/shared/infra/out/messaging"
	"github.com/soat13/fase-1-oficina/internal/shared/messaging"
	helper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/soat13/fase-1-oficina/pkg/observability"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
)

type Container struct {
	SQL               *sql.DB
	DB                *bun.DB
	FiberApp          *fiber.App
	FiberErrorHandler *helper.ErrorHandler
	Validator         *validator.Validate
	EventBus          messaging.Bus
	Metrics           *observability.Metrics
}

func BuildDefault() *Container {
	return Build(nil, nil, nil, nil, nil)
}

func Build(
	bunDB *bun.DB,
	fiberApp *fiber.App,
	fiberErrorHandler *helper.ErrorHandler,
	structValidator *validator.Validate,
	EventBus messaging.Bus,
) *Container {

	if fiberApp == nil {
		fiberApp = newApp()
	}

	if structValidator == nil {
		structValidator = validator.New()
	}

	if fiberErrorHandler == nil {
		fiberErrorHandler = helper.NewErrorHandler(*NewErrorResolver(), structValidator)
	}

	if bunDB == nil {
		bunDB = bun.NewDB(newSQL(), pgdialect.New())
	}

	if EventBus == nil {
		EventBus = infraMessaging.NewInMemoryBus()
	}

	return &Container{
		DB:                bunDB,
		FiberApp:          fiberApp,
		FiberErrorHandler: fiberErrorHandler,
		Validator:         structValidator,
		EventBus:          EventBus,
	}
}

func (c *Container) Close() {
	if c.DB != nil {
		_ = c.DB.Close()
	}
	if c.SQL != nil {
		_ = c.SQL.Close()
	}
}

func newSQL() *sql.DB {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		panic("env PG_DSN not found")
	}

	sqltrace.Register("pgx", &stdlib.Driver{})
	sqlDB, err := sqltrace.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := sqlDB.PingContext(context.Background()); err != nil {
		log.Fatal(err)
	}
	return sqlDB
}

func newApp() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // todo: In production, set specific allowed origins
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD",
		AllowHeaders: "*",
		MaxAge:       3600,
	}))

	app.Use(func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})
	return app
}
