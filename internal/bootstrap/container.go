package bootstrap

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Container struct {
	SQL               *sql.DB
	DB                *bun.DB
	FiberApp          *fiber.App
	FiberErrorHandler *fiberHelper.ErrorHandler
}

func BuildDefault() *Container {
	return Build(nil, nil, nil)
}

func Build(bunDB *bun.DB, fiberApp *fiber.App, fiberErrorHandler *fiberHelper.ErrorHandler) *Container {

	if fiberApp == nil {
		fiberApp = newApp()
	}

	if fiberErrorHandler == nil {
		fiberErrorHandler = fiberHelper.NewErrorHandler(NewErrorResolver())
	}

	if bunDB == nil {
		bunDB = bun.NewDB(newSQL(), pgdialect.New())
	}

	return &Container{
		DB:                bunDB,
		FiberApp:          fiberApp,
		FiberErrorHandler: fiberErrorHandler,
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
	return app
}
