package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
	"github.com/soat13/fase-1-oficina/internal/vehicle/domain"
)

func ErrorHandler(err error, c echo.Context) {
	switch err {
	case domain.ErrInvalidBrand, domain.ErrInvalidModel, domain.ErrInvalidPlate, domain.ErrInvalidYear, domain.ErrInvalidCustomerId:
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

func writeError(ctx *fiber.Ctx, status int, code string, msg string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"code":    code,
		"message": msg,
	})
}
