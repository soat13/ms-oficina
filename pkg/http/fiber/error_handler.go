package fiber

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	errorHelper "github.com/soat13/fase-1-oficina/pkg/error"
)

type ErrorHandler struct {
	ErrorResolver *errorHelper.Resolver
}

func NewErrorHandler(errorResolver *errorHelper.Resolver) *ErrorHandler {

	if errorResolver == nil {
		errorResolver = errorHelper.NewErrorResolver()
	}

	return &ErrorHandler{
		ErrorResolver: errorResolver,
	}
}

func (e *ErrorHandler) Handle(ctx *fiber.Ctx, err error) error {

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    "INVALID_BODY",
			"message": err.Error(),
		})
	}

	if errorInfo, success := e.ErrorResolver.Resolve(err); success {
		return ctx.Status(StatusCodeFromErrorInfo(errorInfo.PrivateCode)).JSON(jsonFromErrorInfo(errorInfo))
	}

	log.Error().
		Err(err).
		Str("path", ctx.Path()).
		Str("method", ctx.Method()).
		Msg("unexpected error on request")

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"code":    "INTERNAL_ERROR",
		"message": "internal error",
	})
}

func StatusCodeFromErrorInfo(status string) int {
	code, err := strconv.Atoi(status)
	if err != nil {
		return fiber.StatusInternalServerError
	}
	return code
}

func jsonFromErrorInfo(errorInfo errorHelper.Info) map[string]string {
	return map[string]string{
		"code":    errorInfo.PublicCode,
		"message": errorInfo.Message(),
	}
}
