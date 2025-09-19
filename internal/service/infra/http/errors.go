package http

import "github.com/gofiber/fiber/v2"

func writeError(ctx *fiber.Ctx, status int, code string, msg string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"code":    code,
		"message": msg,
	})
}
