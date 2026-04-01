package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	authApp "github.com/soat13/fase-1-oficina/internal/auth/application"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	fiberHelper "github.com/soat13/oficina-utils/pkg/http/fiber"
)

type Handler struct {
	authenticate *authApp.AuthenticateUser
	validator    *validator.Validate
	errorHandler *fiberHelper.ErrorHandler
}

func NewHandler(
	authenticate *authApp.AuthenticateUser,
	validator *validator.Validate,
	errorHandler *fiberHelper.ErrorHandler,
) *Handler {
	handler := &Handler{
		authenticate: authenticate,
		validator:    validator,
		errorHandler: errorHandler,
	}

	handler.errorHandler.ErrorResolver.RegisterHTTPUnauthorizedError(authApp.ErrInvalidCredentials)

	return handler
}

func Register(app *fiber.App, h *Handler) {
	app.Post("/auth/login", h.login)
}

type loginBody struct {
	CPF      string `json:"cpf"    validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (h *Handler) login(ctx *fiber.Ctx) error {
	var body loginBody
	if err := ctx.BodyParser(&body); err != nil {
		return h.errorHandler.Handle(ctx, sharedErrors.ErrInvalidJSON)
	}
	if err := h.validator.Struct(body); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	out, err := h.authenticate.Execute(ctx.UserContext(), authApp.AuthenticateInput{
		CPF:      body.CPF,
		Password: body.Password,
	})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	resp := loginResponse{
		AccessToken: out.Token.Value,
		TokenType:   "Bearer",
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
