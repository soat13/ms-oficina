package bootstrap

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	eventbusInfra "github.com/soat13/fase-1-oficina/internal/infra/eventbus"
	"github.com/soat13/fase-1-oficina/internal/ports/event"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Container struct {
	SQL               *sql.DB
	DB                *bun.DB
	FiberApp          *fiber.App
	FiberErrorHandler *fiberHelper.ErrorHandler
	Validator         *validator.Validate
	EventBus          event.Bus
}

func BuildDefault() *Container {
	return Build(nil, nil, nil, nil, nil)
}

func Build(
	bunDB *bun.DB,
	fiberApp *fiber.App,
	fiberErrorHandler *fiberHelper.ErrorHandler,
	structValidator *validator.Validate,
	EventBus event.Bus,
) *Container {

	if fiberApp == nil {
		fiberApp = newApp()
	}

	if structValidator == nil {
		structValidator = validator.New()
	}

	if fiberErrorHandler == nil {
		fiberErrorHandler = fiberHelper.NewErrorHandler(NewErrorResolver(), structValidator)
	}

	if bunDB == nil {
		bunDB = bun.NewDB(newSQL(), pgdialect.New())
	}

	if EventBus == nil {
		EventBus = eventbusInfra.NewInMemoryBus()
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

	sqlDB, err := sql.Open("pgx", dsn)
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
	app.Use(logger.New())
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
