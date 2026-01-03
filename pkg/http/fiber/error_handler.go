package fiber

import (
	"errors"
	"strconv"

	enLocale "github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslator "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	errorHelper "github.com/soat13/fase-1-oficina/pkg/error"
	"github.com/soat13/fase-1-oficina/pkg/maps"
)

type ErrorHandler struct {
	ErrorResolver        errorHelper.Resolver
	validationTranslator ut.Translator
}

func NewErrorHandler(errorResolver errorHelper.Resolver, validator *validator.Validate) *ErrorHandler {

	universalTranslator := ut.New(enLocale.New(), enLocale.New())
	translator, _ := universalTranslator.GetTranslator("enTranslator")
	_ = enTranslator.RegisterDefaultTranslations(validator, translator)

	translator, _ = universalTranslator.GetTranslator("enTranslator")

	return &ErrorHandler{
		ErrorResolver:        errorResolver,
		validationTranslator: translator,
	}
}

func (e *ErrorHandler) Handle(ctx *fiber.Ctx, err error) error {

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    "INVALID_BODY",
			"message": "Body validation failed",
			"errors":  e.getValidationErrors(validationErrors),
		})
	}

	if errorInfo, success := e.ErrorResolver.Resolve(err); success {
		return ctx.Status(statusCodeFromErrorInfo(errorInfo.PrivateCode)).JSON(jsonFromErrorInfo(errorInfo))
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

func statusCodeFromErrorInfo(status string) int {
	code, err := strconv.Atoi(status)
	if err != nil {
		return fiber.StatusInternalServerError
	}
	return code
}

func jsonFromErrorInfo(errorInfo errorHelper.Info) map[string]string {
	return map[string]string{
		"code":    errorInfo.Message(),
		"message": errorInfo.PublicCode,
	}
}

func (e *ErrorHandler) getValidationErrors(err validator.ValidationErrors) []map[string]string {
	return maps.Map(err, func(ve validator.FieldError) map[string]string {
		return map[string]string{
			"field":   ve.Field(),
			"message": ve.Translate(e.validationTranslator),
		}
	})
}
