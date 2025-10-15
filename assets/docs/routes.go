package docs

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
)

var openapiSpec []byte

var swaggerIndex []byte

func Register(app *fiber.App) {
	app.Get("/openapi.yaml", func(c *fiber.Ctx) error {
		c.Type("yaml")
		return c.Send(openapiSpec)
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.Send(swaggerIndex)
	})
}
