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
	"github.com/soat13/oficina-utils/pkg/awsconfig"
	helper "github.com/soat13/oficina-utils/pkg/http/fiber"
	"github.com/soat13/oficina-utils/pkg/messaging"
	sqs "github.com/soat13/oficina-utils/pkg/messaging/sqs"
	"github.com/soat13/oficina-utils/pkg/observability"
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
	Broker            messaging.QueueBroker
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
	broker messaging.QueueBroker,
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

	if broker == nil {
		var err error
		sqsBaseUrl := os.Getenv("SQS_BASE_URL")
		awsConfig := awsconfig.Config{EndpointURL: os.Getenv("AWS_ENDPOINT_URL")}

		broker, err = sqs.NewBroker(context.Background(), awsConfig, sqsBaseUrl)
		if err != nil {
			log.Fatalf("failed to create SQS broker: %v", err)
		}
	}

	return &Container{
		DB:                bunDB,
		FiberApp:          fiberApp,
		FiberErrorHandler: fiberErrorHandler,
		Validator:         structValidator,
		Broker:            broker,
	}
}

func (c *Container) StartConsumers(ctx context.Context) {
	c.Broker.Listen(ctx)
}

func (c *Container) Publisher() messaging.QueueSender {
	return c.Broker
}

func (c *Container) Subscribe(topic string, handler messaging.Handler) {
	c.Broker.Subscribe(topic, handler)
}

func (c *Container) Close() {
	if c.Broker != nil {
		c.Broker.Stop()
	}
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
